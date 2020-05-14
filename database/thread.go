package database

import (
	"context"

	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/logger"
)

// GetThreadByURI get thread by uri
func (d *Database) GetThreadByURI(ctx context.Context, uri string) (isso.Thread, error) {
	logger.Debug("uri %s", uri)
	var thread isso.Thread
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	err := d.DB.QueryRowContext(ctx, d.statement["thread_get_by_uri"], uri).Scan(&thread.ID, &thread.URI, &thread.Title)
	if err != nil {
		return thread, wraperror(err)
	}
	return thread, nil
}

// GetThreadByID get thread by id
func (d *Database) GetThreadByID(ctx context.Context, id int64) (isso.Thread, error) {
	logger.Debug("id %d", id)
	var thread isso.Thread
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	err := d.DB.QueryRowContext(ctx, d.statement["thread_get_by_id"], id).Scan(&thread.ID, &thread.URI, &thread.Title)
	if err != nil {
		return thread, wraperror(err)
	}
	return thread, nil
}

// NewThread new a thread
func (d *Database) NewThread(ctx context.Context, uri string, title string) (isso.Thread, error) {
	logger.Debug("create thread %s %s", uri, title)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	if title == "" || uri == "" {
		return isso.Thread{}, wraperror(isso.ErrInvalidParam)
	}

	var rowsaffected, lastinsertid int64
	err := d.execstmt(ctx, &rowsaffected, &lastinsertid, d.statement["thread_new"], uri, title)
	if err != nil {
		return isso.Thread{}, wraperror(err)
	}
	if rowsaffected != 1 {
		return isso.Thread{}, wraperror(isso.ErrNotExpectAmount)
	}
	return isso.Thread{ID: lastinsertid, URI: uri, Title: title}, nil
}
