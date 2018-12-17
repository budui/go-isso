package db

import (
	"database/sql"
	"log"

	"github.com/RayHY/go-isso/internal/pkg/conf"

	// can be easily replaced by mysql etc.
	_ "github.com/mattn/go-sqlite3"
)

var (
	createComments = `
	CREATE TABLE IF NOT EXISTS comments (
		tid REFERENCES threads(id), 
		id INTEGER PRIMARY KEY, 
		parent INTEGER,
		created FLOAT NOT NULL,
		modified FLOAT,
		mode INTEGER,
		remote_addr VARCHAR,
		text VARCHAR,
		author VARCHAR,
		email VARCHAR,
		website VARCHAR,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		voters BLOB NOT NULL,
		notification INTEGER DEFAULT 0
	);
	`
	createThreads = `
	CREATE TABLE IF NOT EXISTS threads (
		id INTEGER PRIMARY KEY,
		uri VARCHAR(256) UNIQUE,
		title VARCHAR(256)
	);
	`

	createPreferences = `
	CREATE TABLE IF NOT EXISTS preferences (
		key VARCHAR PRIMARY KEY, 
		value VARCHAR
	);
	`
	
	createAutomateDorphanRemoval = `
	CREATE TRIGGER IF NOT EXISTS remove_stale_threads
    AFTER DELETE ON comments
    BEGIN
    	DELETE FROM threads WHERE id NOT IN (SELECT tid FROM comments);
    END
	`

	// MAYBE replace it with something like pragma_table_info('comments')?
	convertOldIssoDatabase = `
		ALTER TABLE comments ADD COLUMN notification INTEGER DEFAULT 0;
	`
)

// Accessor defines all usual access ops avail.
type Accessor interface {
	commentAccessor
	threadsAccessor
	Close() error
}

type database struct {
	*sql.DB
	guard conf.Guard
}

// NewAccessor generate an new DB worker
// 1. Open database.
// 2. Ping database.
// 3. Create Table comment, thread, preference if they do not exist.
// 4. Add field `notification` when use old isso database.
func NewAccessor(path string, guard conf.Guard) (Accessor, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// if need to add another database support, just use another sql slice.
	Sqlite3createSQL := []string{createComments, createPreferences, createThreads, createAutomateDorphanRemoval}

	for _, InitSQL := range Sqlite3createSQL {
		_, err := db.Exec(InitSQL)
		if err != nil {
			return nil, err
		}
	}
	// Add field `notification` when use old isso database.
	// failed when exist `notification`.
	db.Exec(convertOldIssoDatabase)

	return &database{db, guard}, nil
}

func (db *database) Close() error {
	return db.Close()
}
