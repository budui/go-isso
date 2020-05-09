package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/guregu/null.v4"
	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/logger"
)

// Database handles all operations related to the database.
type Database struct {
	*sql.DB
	statement map[string]string
	timeout   time.Duration
}

// ErrNotExpectRow is returned by database method when affected row is not equal as expect.
var ErrNotExpectRow = errors.New("database: affected row is not equal as expect")

// New return a *Database
func New(path string, timeout time.Duration) (*Database, error) {
	databaseType := "sqlite3"
	if path == "" {
		path = ":memory:"
	}

	db, err := sql.Open(databaseType, path)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	_, err = db.Exec(presetSQL[databaseType]["create"])
	if err != nil {
		return nil, err
	}
	// Add field `notification` when use old isso database.
	// failed when exist `notification`.
	// so just IGNORE error.
	_, _ = db.Exec(presetSQL[databaseType]["migrate_add_notification"])
	logger.Debug("create database instance at %s", path)
	return &Database{db, presetSQL[databaseType], timeout}, nil
}

func (d *Database) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d.timeout)
}

type nullComment struct {
	TID          int64
	ID           int64
	Parent       null.Int
	Created      float64
	Modified     null.Float
	Mode         int
	RemoteAddr   string
	Text         string
	Author       string
	Email        null.String
	Website      null.String
	Likes        int
	Dislikes     int
	Voters       []byte
	Notification int
}

func (nc nullComment) ToComment() isso.Comment {
	c := isso.Comment{
		ID:           nc.ID,
		Parent:       &nc.Parent.Int64,
		Created:      nc.Created,
		Modified:     &nc.Modified.Float64,
		Mode:         nc.Mode,
		Text:         nc.Text,
		Author:       nc.Author,
		Email:        &nc.Email.String,
		Website:      &nc.Website.String,
		Likes:        nc.Likes,
		Dislikes:     nc.Dislikes,
		Notification: nc.Notification,
		RemoteAddr:   nc.RemoteAddr,
	}
	if !nc.Parent.Valid {
		c.Parent = nil
	}
	if !nc.Modified.Valid {
		c.Modified = nil
	}
	if !nc.Email.Valid {
		c.Email = nil
	}
	if !nc.Website.Valid {
		c.Website = nil
	}
	return c
}

func newNullComment(c isso.Comment, threadID int64, remoteAddr string) nullComment {
	// TODO: generated votes use Bloomfilter
	votes := []byte{}
	return nullComment{
		TID:          threadID,
		ID:           c.ID,
		Parent:       null.IntFromPtr(c.Parent),
		Created:      float64(time.Now().UnixNano()) / float64(1e9),
		Modified:     null.NewFloat(0, false),
		Mode:         c.Mode,
		RemoteAddr:   remoteAddr,
		Text:         c.Text,
		Author:       c.Author,
		Email:        null.StringFromPtr(c.Email),
		Website:      null.StringFromPtr(c.Website),
		Likes:        c.Likes,
		Dislikes:     c.Dislikes,
		Voters:       votes,
		Notification: c.Notification,
	}
}
