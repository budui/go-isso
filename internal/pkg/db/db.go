package db

import (
	"database/sql"
	"errors"
	"github.com/RayHY/go-isso/internal/pkg/conf"
	"gopkg.in/guregu/null.v3"
	"log"
	// can be easily replaced by mysql etc.
	"github.com/RayHY/go-isso/internal/app/isso/util"
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

	// MAYBE replace it with something like pragma_table_info('comments')?
	convertOldIssoDatabase = `
		ALTER TABLE comments ADD COLUMN notification INTEGER DEFAULT 0;
	`
)

// Accessor defines all usual access ops avail.
type Accessor interface {
	Close() error

	// CURD stuff

	// Add new comment to DB and return a complete Comment.
	Add(uri string, c Comment) (Comment, error)
	// Update comment `id` with values from `data`
	Update(id int64, c Comment) (Comment, error)
	// Search for comment `id` and return a mapping of `fields` and values.
	Get(id int64) (Comment, error)
	// Return comment count for main thread and all reply threads for one url.
	CountReply(uri string, mode int, after float64) (map[null.Int]int64, error)
	// Return comments for `uri` with `mode`.
	Fetch(uri string, mode int, after float64, parent null.Int, orderBy string, isASC bool, limit null.Int) ([]Comment, error)
}

type database struct {
	*sql.DB
	guard conf.Guard
}

// NewWorker generate an new DB worker
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
	Sqlite3createSQL := []string{createComments, createPreferences, createThreads}

	for _, InitSql := range Sqlite3createSQL {
		_, err := db.Exec(InitSql)
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

// CountReply return comment count for main thread and all reply threads for one url.
func (db *database) CountReply(uri string, mode int, after float64) (map[null.Int]int64, error) {
	countSQL := `SELECT comments.parent,count(*)
            FROM comments INNER JOIN threads ON
            	threads.uri=? AND comments.tid=threads.id AND
               	(? | comments.mode = ?) AND
               	comments.created > ?
			GROUP BY comments.parent
			`
	rows, err := db.Query(countSQL, uri, mode, mode, after)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	countResult := make(map[null.Int]int64)
	var parent null.Int
	var count int64
	for rows.Next() {
		err := rows.Scan(&parent, &count)
		if err != nil {
			return nil, err
		}
		countResult[parent] = count
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return countResult, nil
}

// Fetch Return comments for `uri` with `mode`.
// parent: <0 for any.
// opderBy can be select from ['id', 'created', 'modified', 'likes', 'dislikes']
// isasc: true for ASC, false for DESC
func (db *database) Fetch(uri string, mode int, after float64, parent null.Int, orderBy string, isASC bool, limit null.Int) ([]Comment, error) {
	statement := `SELECT comments.* 
			FROM comments INNER JOIN threads ON
            	threads.uri=? AND comments.tid=threads.id AND (? | comments.mode) = ?
				AND comments.created>?
	`

	switch {
	case !parent.Valid:
		statement += `AND comments.parent IS NULL`
	case parent.ValueOrZero() >= 0:
		statement += `AND comments.parent=?`
	}

	switch orderBy {
	case "id", "created", "modified", "likes", "dislikes":
		break
	default:
		orderBy = "id"
	}
	statement += ` ORDER BY `
	statement += orderBy

	if !isASC {
		statement += ` DESC `
	}

	var rows *sql.Rows
	var err error
	switch {
	// (top level comments | all comments) without limit
	case (!parent.Valid || parent.Int64 < 0) && !limit.Valid:
		rows, err = db.Query(statement, uri, mode, mode, after)
		// (top level comments | all comments) with limit
	case (!parent.Valid || parent.Int64 < 0) && limit.Int64 > 0:
		statement += ` LIMIT ?`
		rows, err = db.Query(statement, uri, mode, mode, after, limit)
		// specific parent comments with limit
	case limit.Valid && limit.Int64 > 0:
		statement += ` LIMIT ?`
		rows, err = db.Query(statement, uri, mode, mode, after, parent, limit)
		// specific parent comments without limit
	case !limit.Valid:
		rows, err = db.Query(statement, uri, mode, mode, after, parent)
	case limit.Valid && limit.Int64 < 0:
		return nil, errors.New("db.Fetch: invalid limit")
	case limit.Valid && limit.Int64 == 0:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment

		err := rows.Scan(&c.tid, &c.ID, &c.Parent, &c.Created, &c.Modified,
			&c.Mode, &c.remoteAddr, &c.Text, &c.Author, &c.email, &c.Website,
			&c.Likes, &c.Dislikes, &c.voters, &c.notification,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Add new comment to DB and return a compelete Comment
func (db *database) Add(uri string, c Comment) (Comment, error) {
	if c.Parent.Valid {
		parent, err := db.Get(c.Parent.Int64)
		if err != nil {
			return Comment{}, err
		}
		c.Parent = parent.Parent
	}

	c.voters = util.GenBloomfilterfunc(c.remoteAddr)

	stmt, err := db.Prepare(`
	INSERT INTO comments (
        	tid, parent, created, modified, mode, remote_addr,
			text, author, email, website, voters, notification
		)
    SELECT threads.id, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			FROM threads WHERE threads.uri = ?;`)
	if err != nil {
		return Comment{}, err
	}
	_, err = stmt.Exec(c.Parent, c.Created, c.Modified, c.Mode, c.remoteAddr,
		c.Text, c.Author, c.email, c.Website, c.voters, c.notification, uri)
	if err != nil {
		return Comment{}, err
	}

	err = db.QueryRow(`SELECT c.* FROM comments AS c INNER JOIN threads ON threads.uri = ? ORDER BY c.id DESC LIMIT 1`, uri).Scan(
		&c.tid, &c.ID, &c.Parent, &c.Created, &c.Modified,
		&c.Mode, &c.remoteAddr, &c.Text, &c.Author, &c.email, &c.Website,
		&c.Likes, &c.Dislikes, &c.voters, &c.notification,
	)
	if err != nil {
		return Comment{}, err
	}
	return c, nil
}

// Get : Search for comment :param:`id` and return (Comment, nil) or (nil, error)
func (db *database) Get(id int64) (Comment, error) {
	var c Comment
	err := db.QueryRow(`SELECT * FROM comments WHERE id=?`, id).Scan(
		&c.tid, &c.ID, &c.Parent, &c.Created, &c.Modified,
		&c.Mode, &c.remoteAddr, &c.Text, &c.Author, &c.email, &c.Website,
		&c.Likes, &c.Dislikes, &c.voters, &c.notification,
	)
	if err != nil {
		return Comment{}, err
	}
	return c, err
}

func (db *database) Update(id int64, c Comment) (Comment, error) {
	nc := Comment{}
	return nc, nil
}
