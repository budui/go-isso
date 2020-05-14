package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"wrong.wang/x/go-isso/extract"
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
func (d *Database) NewThread(ctx context.Context, uri string, title string, Host string) (isso.Thread, error) {
	logger.Debug("create thread %s host: %s", uri, Host)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	if title == "" {
		url := Host + uri
		logger.Debug("database: fetch %s", url)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return isso.Thread{}, wraperror(fmt.Errorf("get title failed: do request failed %v, %w", url, err))
		}
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			return isso.Thread{}, wraperror(fmt.Errorf("get title failed: %s, %w", url, err))
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return isso.Thread{}, wraperror(fmt.Errorf("get title failed: %s, %d", url, resp.StatusCode))
		}
		title, uri, err = extract.TitleAndThreadURI(resp.Body, "Untitled", uri)
		if err != nil {
			return isso.Thread{}, wraperror(fmt.Errorf("get title failed: %s, %w", url, err))
		}
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
