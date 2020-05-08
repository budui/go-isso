package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	err := d.DB.QueryRow(d.statement["comment_get_by_id"], id).Scan(
		&nc.TID, &nc.ID, &nc.Parent, &nc.Created, &nc.Modified, &nc.Mode,
		&nc.RemoteAddr, &nc.Text, &nc.Author, &nc.Email, &nc.Website, &nc.Likes,
		&nc.Dislikes, &nc.Voters, &nc.Notification,
	)
	if err != nil {
		return nc, err
	}
	return nc, nil
}
