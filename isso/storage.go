package isso

import "context"

// Storage handles all operations related to the database.
type Storage interface {
	ThreadStorage
	CommentStorage
	PreferenceStorage
}

// ThreadStorage handles all operations related to Thread and the database.
type ThreadStorage interface {
	GetThreadByURI(ctx context.Context, uri string) (Thread, error)
	GetThreadByID(ctx context.Context, id int64) (Thread, error)
	NewThread(ctx context.Context, uri string, title string, Host string) (Thread, error)
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
