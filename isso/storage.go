package isso

// Storage handles all operations related to the database.
type Storage interface {
	ThreadStorage
	CommentStorage
	PreferenceStorage
}

// ThreadStorage handles all operations related to Thread and the database.
type ThreadStorage interface {
	ContainsThread(uri string) bool
	GetThreadByURI(uri string) (Thread, error)
	GetThreadByID(id int) (Thread, error)
	NewThread(uri string, title string, Host string) (Thread, error)
}

// CommentStorage handles all operations related to Comment and the database.
type CommentStorage interface {
	IsApprovedAuthor(email string) bool
	NewComment(c Comment, threadID int64, remoteAddr string) (Comment, error)
}

// PreferenceStorage handles all operations related to Preference and the database.
type PreferenceStorage interface {
	GetPreference(key string) (string, error)
	SetPreference(key string, value string) error
}

