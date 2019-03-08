package db

import "gopkg.in/guregu/null.v3"

// Accessor defines all usual access ops avail.
type Accessor interface {
	commentAccessor
	threadsAccessor
	preferenceAccessor
	Close() error
}

// commentAccessor defines all usual access ops avail for comment.
type commentAccessor interface {
	// CURD stuff

	// Add new comment to DB and return a complete Comment.
	Add(uri string, c Comment) (Comment, error)
	// Delete a comment.
	Delete(id int64) (Comment, error)
	// Update comment `id` with values from `data`
	Update(id int64, text string, author, website null.String) (Comment, error)
	// Search for comment `id` and return a mapping of `fields` and values.
	Get(id int64) (Comment, error)
	// Return comment count for main thread and all reply threads for one url.
	CountReply(uri string, mode int, after float64) (map[null.Int]int64, error)
	// Return comments for `uri` with `mode`.
	Fetch(uri string, mode int, after float64, parent null.Int, orderBy string, isASC bool, limit null.Int) ([]Comment, error)
}

// preferenceAccessor defines all usual access ops avail for comment.
type preferenceAccessor interface {
	GetPreference(key string) (null.String, error)
	SetPreference(key, value string) error
}

// threadsAccessor defines all usual access ops avail for threads
type threadsAccessor interface {
	Contain(uri string) (bool, error)
	GetThreadWithID(ID int64) (Thread, error)
	GetThreadWithURI(uri string) (Thread, error)
	NewThread(uri string, title null.String) (Thread, error)
}
