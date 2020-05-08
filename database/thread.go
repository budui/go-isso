package database

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	"wrong.wang/x/go-isso/extract"
	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/logger"
)

// GetThreadByURI get thread by uri
func (d *Database) GetThreadByURI(ctx context.Context, uri string) (isso.Thread, error) {
	logger.Debug("database: get thread by %s", uri)
	var thread isso.Thread
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	err := d.DB.QueryRowContext(ctx, d.statement["thread_get_by_uri"], uri).Scan(&thread.ID, &thread.URI, &thread.Title)
	if err != nil {
		return thread, fmt.Errorf("GetThreadByURI failed. %w", err)
	}
	return thread, nil
}

// GetThreadByID get thread by id
func (d *Database) GetThreadByID(ctx context.Context, id int64) (isso.Thread, error) {
	logger.Debug("database: get thread by %d", id)
	var thread isso.Thread
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()
	err := d.DB.QueryRowContext(ctx, d.statement["thread_get_by_id"], id).Scan(&thread.ID, &thread.URI, &thread.Title)
	if err != nil {
		return thread, fmt.Errorf("GetThreadByURI failed. %w", err)
	}
	return thread, nil
}

// NewThread new a thread
func (d *Database) NewThread(ctx context.Context, uri string, title string, Host string) (isso.Thread, error) {
	logger.Debug("database: create thread %s host: %s", uri, Host)
	ctx, cancel := d.withTimeout(ctx)
	defer cancel()

	if title == "" {
		url := path.Join(Host, uri)
		logger.Debug("database: fetch %s", url)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return isso.Thread{}, fmt.Errorf("NewThread: get title failed. failed to load page %s, %w", url, err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			return isso.Thread{}, fmt.Errorf("NewThread: get title failed. failed to load page %s, %w", url, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return isso.Thread{}, fmt.Errorf("NewThread: get title failed. can't load page %s, code %d", url, resp.StatusCode)
		}
		title, uri, err = extract.TitleAndThreadURI(resp.Body, "Untitled", uri)
		if err != nil {
			return isso.Thread{}, fmt.Errorf("NewThread: get title failed. can't extract title, %w", err)
		}
	}

	result, err := d.DB.ExecContext(ctx, d.statement["thread_new"], uri, title)
	if err != nil {
		return isso.Thread{}, fmt.Errorf("NewThread failed. %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return isso.Thread{}, fmt.Errorf("NewThread failed. %w", err)
	}
	return isso.Thread{ID: id, URI: uri, Title: title}, nil
}
