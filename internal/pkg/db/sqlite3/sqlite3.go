package sqlite3

import (
	"database/sql"
	"github.com/RayHY/go-isso/internal/pkg/conf"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// sqlite3CreateSQLs represents a list of queries for initial table creation in sqlite
var sqlite3CreateSQLs = []string{
	`CREATE TABLE IF NOT EXISTS comments (
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
	);`,
	`CREATE TABLE IF NOT EXISTS preferences (
		key VARCHAR PRIMARY KEY, 
		value VARCHAR
	);`,
	`CREATE TABLE IF NOT EXISTS threads (
		id INTEGER PRIMARY KEY,
		uri VARCHAR(256) UNIQUE,
		title VARCHAR(256)
	);
	`, `
	CREATE TRIGGER IF NOT EXISTS remove_stale_threads
    AFTER DELETE ON comments
    BEGIN
    	DELETE FROM threads WHERE id NOT IN (SELECT tid FROM comments);
    END
	`,
}

// MAYBE replace it with something like pragma_table_info('comments')?
var convertOldIssoDatabase = `ALTER TABLE comments ADD COLUMN notification INTEGER DEFAULT 0;`

func CreateDatabase(sqliteConf conf.Sqlite3) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sqliteConf.Path)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	for _, InitSQL := range sqlite3CreateSQLs {
		_, err := db.Exec(InitSQL)
		if err != nil {
			return nil, err
		}
	}
	// Add field `notification` when use old isso database.
	// failed when exist `notification`.
	// so just IGNORE error.
	_, _ = db.Exec(convertOldIssoDatabase)
	return db, nil
}
