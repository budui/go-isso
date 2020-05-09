package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gopkg.in/guregu/null.v4"
	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/logger"
)

// IsApprovedAuthor check if email has approved in 6 month
func (d *Database) IsApprovedAuthor(ctx context.Context, email string) bool {
	logger.Debug("database: check email %s", email)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	if email == "" {
		return false
	}
	var flag int64
	err := d.DB.QueryRowContext(ctx, d.statement["comment_is_previously_approved_author"], email).Scan(&flag)
	return (err == nil) && (flag == 1)
}

// NewComment add comment into database
func (d *Database) NewComment(ctx context.Context, c isso.Comment, threadID int64, remoteAddr string) (isso.Comment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	logger.Debug("database: create comment at %d", threadID)
	if c.Parent != nil {
		parent, err := d.getComment(ctx, *c.Parent)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return isso.Comment{}, errors.New("NewComment: can't find specify parent")
			}
			return isso.Comment{}, errors.New("NewComment: find specify parent failed")
		}
		if parent.TID != threadID {
			return isso.Comment{}, errors.New("NewComment: not same tid with parent's tid")
		}
		if parent.Parent.Valid {
			c.Parent = &parent.Parent.Int64
		}
	}

	nc := newNullComment(c, threadID, remoteAddr)

	result, err := d.DB.ExecContext(ctx, d.statement["comment_new"], nc.TID, nc.Parent, nc.Created,
		nc.Modified, nc.Mode, nc.RemoteAddr, nc.Text, nc.Author, nc.Email, nc.Website, nc.Voters, nc.Notification,
	)
	if err != nil {
		return isso.Comment{}, fmt.Errorf("NewComment: insert failed - %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return isso.Comment{}, fmt.Errorf("NewComment: insert failed - %w", err)
	}
	comment, err := d.GetComment(ctx, id)
	if err != nil {
		return isso.Comment{}, fmt.Errorf("NewComment: insert failed - %w", err)
	}
	return comment, nil
}

// GetComment get comment by ID
func (d *Database) GetComment(ctx context.Context, id int64) (isso.Comment, error) {
	logger.Debug("database: get comment %d", id)
	nc, err := d.getComment(ctx, id)
	if err != nil {
		return isso.Comment{}, fmt.Errorf("GetComment: failed %w", err)
	}
	return nc.ToComment(), nil
}

func (d *Database) getComment(ctx context.Context, id int64) (nullComment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	var nc nullComment
	err := d.DB.QueryRowContext(ctx, d.statement["comment_get_by_id"], id).Scan(
		&nc.TID, &nc.ID, &nc.Parent, &nc.Created, &nc.Modified, &nc.Mode,
		&nc.RemoteAddr, &nc.Text, &nc.Author, &nc.Email, &nc.Website, &nc.Likes,
		&nc.Dislikes, &nc.Voters, &nc.Notification,
	)
	if err != nil {
		return nc, err
	}
	return nc, nil
}

// CountReply return comment count for main thread's comment and all reply threads for one uri.
// 0 mean null parent
func (d *Database) CountReply(ctx context.Context, uri string, mode int, after float64) (map[int64]int64, error) {
	logger.Debug("database: CountReplyPerComment %s", uri)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	var topCommentsCount int64
	var counts map[int64]int64

	rows, err := d.DB.QueryContext(ctx, d.statement["comment_count_reply"], uri, mode, mode, after)
	if err != nil {
		return nil, fmt.Errorf("CountReplyPerComment failed %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p null.Int
		var c int64
		err := rows.Scan(&p, &c)
		if err != nil {
			return nil, fmt.Errorf("CountReplyPerComment failed %w", err)
		}
		if p.Valid {
			counts[p.Int64] = c
		} else {
			topCommentsCount = c
		}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("CountReplyPerComment failed %w", err)
	}
	counts[0] = topCommentsCount
	return counts, nil
}

// FetchCommentsByURI fetch comments related uri with a lot of param
func (d *Database) FetchCommentsByURI(ctx context.Context, uri string, parent int64, mode int, orderBy string, asc bool) (map[int64][]isso.Comment, error) {
	logger.Debug("database: FetchComments %s", uri)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	switch orderBy {
	case "id", "created", "modified", "likes", "dislikes":
	default:
		orderBy = "id"
	}

	desc := ""
	if !asc {
		desc += ` DESC `
	}

	condition := fmt.Sprintf(" ORDER BY %s %s", orderBy, desc)

	var rows *sql.Rows
	var err error
	switch {
	case parent < 0:
		stmt := d.statement["comment_fetch_by_uri"] + condition
		rows, err = d.DB.QueryContext(ctx, stmt, uri, mode, mode)
	case parent == 0:
		stmt := d.statement["comment_fetch_by_uri"] + ` AND comments.parent IS NULL ` + condition
		rows, err = d.DB.QueryContext(ctx, stmt, uri, mode, mode)
	case parent > 0:
		stmt := d.statement["comment_fetch_by_uri"] + ` AND comments.parent=? ` + condition
		rows, err = d.DB.QueryContext(ctx, stmt, uri, mode, mode, parent)
	}

	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("FetchCommentsByURI failed %w", err)
	}

	var commentsbyparent map[int64][]isso.Comment

	for rows.Next() {
		var nc nullComment

		err := rows.Scan(
			&nc.TID, &nc.ID, &nc.Parent, &nc.Created, &nc.Modified, &nc.Mode,
			&nc.RemoteAddr, &nc.Text, &nc.Author, &nc.Email, &nc.Website, &nc.Likes,
			&nc.Dislikes, &nc.Voters, &nc.Notification,
		)
		if err != nil {
			return nil, fmt.Errorf("FetchCommentsByURI failed %w", err)
		}
		if nc.Parent.Valid {
			commentsbyparent[nc.Parent.Int64] = append(commentsbyparent[nc.Parent.Int64], nc.ToComment())
		} else {
			commentsbyparent[0] = append(commentsbyparent[0], nc.ToComment())
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("FetchCommentsByURI failed %w", err)
	}
	return commentsbyparent, nil
}
