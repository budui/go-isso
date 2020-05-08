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
}

// PreferenceStorage handles all operations related to Preference and the database.
type PreferenceStorage interface {
	GetPreference(key string) (string, error)
	SetPreference(key string, value string) error
}
