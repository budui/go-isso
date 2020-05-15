package isso

import (
	"context"
	"errors"
)

// predictable error
var (
	// ErrStorageNotFound is returned by Storage when no result can be found
	ErrStorageNotFound = errors.New("storage: no result found")
	// ErrNotExpectAmount is returned by databStoragease method when affected amount is not equal as expect.
	ErrNotExpectAmount = errors.New("storage: affected amount is not equal as expect")
	// ErrInvalidParam is returned by databStoragease method when input param check failed.
	ErrInvalidParam = errors.New("storage: handler input is not valid")
)

// mode for comment's mode. comment mode CAN NOT be set to modePublic.
const (
	//modeAccepted means The comment was accepted by the server and is published.
	// 001
	ModeAccepted = 1
	//modeModeration means: The comment was accepted by the server but awaits moderation.
	// 010
	ModeModeration = 2
	//modeDeleted means deleted, but referenced: The comment was deleted on the server but is still referenced by replies.
	// 100
	ModeDeleted = 4
	//modePublic means The comment is public. its replies can be counted and show.
	// 101
	//modePublic include modeAccepted, modeDeleted
	//modePublic CAN NOT be used by comment.
	//It is the shortcut for select comments who are in modeAccepted or modeDeleted.
	ModePublic = 5
)

// Storage handles all operations related to the database.
type Storage interface {
	ThreadStorage
	CommentStorage
	PreferenceStorage
	NewCommentGuard(ctx context.Context, c Comment, uri string,
		ratelimit int, directreply int, replytoself bool, maxage int) (bool, string)
}

// ThreadStorage handles all operations related to Thread and the database.
type ThreadStorage interface {
	GetThreadByURI(ctx context.Context, uri string) (Thread, error)
	GetThreadByID(ctx context.Context, id int64) (Thread, error)
	NewThread(ctx context.Context, uri string, title string) (Thread, error)
}

// CommentStorage handles all operations related to Comment and the database.
type CommentStorage interface {
	IsApprovedAuthor(ctx context.Context, email string) bool
	NewComment(ctx context.Context, c Comment, threadID int64, remoteAddr string) (Comment, error)
	GetComment(ctx context.Context, id int64) (Comment, error)
	// CountReply return parent-count map, 0 mean null `parent`
	CountReply(ctx context.Context, uri string, mode int, after float64) (map[int64]int64, error)
	FetchCommentsByURI(ctx context.Context, uri string, parent int64, mode int, orderBy string, asc bool) (map[int64][]Comment, error)
	CountComment(ctx context.Context, uris []string) (map[string]int64, error)
}

// PreferenceStorage handles all operations related to Preference and the database.
type PreferenceStorage interface {
	GetPreference(key string) (string, error)
	SetPreference(key string, value string) error
}
