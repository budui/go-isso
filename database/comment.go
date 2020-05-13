package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/logger"
)

// IsApprovedAuthor check if email has approved in 6 month
func (d *Database) IsApprovedAuthor(ctx context.Context, email string) bool {
	logger.Debug("email %s", email)
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
	logger.Debug("create %s 's comment at %d", c.Author, threadID)
	if c.Parent != nil {
		parent, err := d.getComment(ctx, *c.Parent)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return isso.Comment{}, wraperror(err)
			}
			return isso.Comment{}, wraperror(err)
		}
		if parent.TID != threadID {
			return isso.Comment{}, wraperror(err)
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
		return isso.Comment{}, wraperror(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return isso.Comment{}, wraperror(err)
	}
	comment, err := d.GetComment(ctx, id)
	if err != nil {
		return isso.Comment{}, wraperror(err)
	}
	return comment, nil
}

// GetComment get comment by ID
func (d *Database) GetComment(ctx context.Context, id int64) (isso.Comment, error) {
	logger.Debug("get comment %d", id)
	nc, err := d.getComment(ctx, id)
	if err != nil {
		return isso.Comment{}, wraperror(err)
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
	logger.Debug("uri: %s", uri)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	counts := map[int64]int64{}

	rows, err := d.DB.QueryContext(ctx, d.statement["comment_count_reply"], uri, mode, mode, after)
	if err != nil {
		return nil, wraperror(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p null.Int
		var c int64
		err := rows.Scan(&p, &c)
		if err != nil {
			return nil, wraperror(err)
		}
		if p.Valid {
			counts[p.Int64] = c
		} else {
			counts[0] = c
		}
	}
	if rows.Err() != nil {
		return nil, wraperror(err)
	}
	return counts, nil
}

// FetchCommentsByURI fetch comments related uri with a lot of param
func (d *Database) FetchCommentsByURI(ctx context.Context, uri string, parent int64, mode int, orderBy string, asc bool) (map[int64][]isso.Comment, error) {
	logger.Debug("uri: %s", uri)
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
		return nil, wraperror(err)
	}

	commentsbyparent := map[int64][]isso.Comment{}

	for rows.Next() {
		var nc nullComment

		err := rows.Scan(
			&nc.TID, &nc.ID, &nc.Parent, &nc.Created, &nc.Modified, &nc.Mode,
			&nc.RemoteAddr, &nc.Text, &nc.Author, &nc.Email, &nc.Website, &nc.Likes,
			&nc.Dislikes, &nc.Voters, &nc.Notification,
		)
		if err != nil {
			return nil, wraperror(err)
		}
		if nc.Parent.Valid {
			commentsbyparent[nc.Parent.Int64] = append(commentsbyparent[nc.Parent.Int64], nc.ToComment())
		} else {
			commentsbyparent[0] = append(commentsbyparent[0], nc.ToComment())
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, wraperror(err)
	}
	return commentsbyparent, nil
}

// CountComment count comment per thread
func (d *Database) CountComment(ctx context.Context, uris []string) (map[string]int64, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	logger.Debug("uris: %v", uris)
	commentByURI := map[string]int64{}
	if len(uris) == 0 {
		return commentByURI, nil
	}
	rows, err := d.DB.QueryContext(ctx, d.statement["comment_count"])
	defer rows.Close()

	if err != nil {
		return nil, wraperror(err)
	}

	for rows.Next() {
		var uri string
		var count int64
		err := rows.Scan(&uri, &count)
		if err != nil {
			return nil, wraperror(err)
		}
		commentByURI[uri] = count
	}

	err = rows.Err()
	if err != nil {
		return nil, wraperror(err)
	}

	uriMap := map[string]int64{}
	for _, uri := range uris {
		uriMap[uri] = commentByURI[uri]
	}
	return uriMap, nil
}

// ActivateComment Activate comment id if pending
func (d *Database) ActivateComment(ctx context.Context, id int64) error {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	logger.Debug("id: %d", id)

	var rowsaffected int64
	err := d.execstmt(ctx, &rowsaffected, nil, d.statement["comment_activate"], id)
	if err != nil {
		return wraperror(err)
	}
	if rowsaffected != 1 {
		return wraperror(isso.ErrNotExpectAmount)
	}
	return nil
}

// EditComment edit comment
func (d *Database) EditComment(ctx context.Context, c isso.Comment) (isso.Comment, error) {
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	logger.Debug("edit %s 's comment", c.Author)

	var rowsaffected int64
	modified := float64(time.Now().UnixNano()) / float64(1e9)
	err := d.execstmt(ctx, &rowsaffected, nil, d.statement["comment_edit"],
		c.Text, c.Author, null.StringFromPtr(c.Website), modified, c.ID)
	if err != nil {
		return isso.Comment{}, wraperror(err)
	}
	if rowsaffected != 1 {
		return isso.Comment{}, wraperror(isso.ErrNotExpectAmount)
	}

	comment, err := d.GetComment(ctx, c.ID)
	if err != nil {
		return isso.Comment{}, wraperror(err)
	}
	return comment, nil
}
