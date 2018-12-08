package db

import (
	"database/sql"
	"errors"
	"github.com/RayHY/go-isso/internal/app/isso/util"

	"gopkg.in/guregu/null.v3"

	log "github.com/sirupsen/logrus"
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
)

// Worker have all database-related methods.
type Worker interface {
	// Prepare prepare all database-related things for isso
	PrepareToWork() error
	CountReply(uri string, mode int, after float64) (map[null.Int]int64, error)
	Fetch(uri string, mode int, after float64, parent null.Int, orderBy string, isASC bool, limit null.Int) ([]Comment, error)
	Add(uri string, c Comment) (Comment, error)
	Get(id int64) (Comment, error)
}

// Guard holds all config for database limit
type Guard struct {
	IsAlive                        bool
	RateLimit                      int
	ReplyLimit                     int
	CanReplyToSelfWhenCanStillEdit bool
	NeedAuthor                     bool
	NeedEmail                      bool
}

type database struct {
	*sql.DB
	guard Guard
}

// NewWorker generate an new DB worker
func NewWorker(path string, guard Guard) Worker {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	return &database{db, guard}
}

// PrepareToWork prepare all database-related things for go-isso
func (db *database) PrepareToWork() error {

	if err := db.Ping(); err != nil {
		return err
	}

	// if need to add another database support, just use another sql slice.
	Sqlite3createSQL := []string{createComments, createPreferences, createThreads}

	for _, sql := range Sqlite3createSQL {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// CountReply return comment count for main thread and all reply threads for one url.
func (db *database) CountReply(uri string, mode int, after float64) (map[null.Int]int64, error) {
	sql := `SELECT comments.parent,count(*)
            FROM comments INNER JOIN threads ON
            	threads.uri=? AND comments.tid=threads.id AND
               	(? | comments.mode = ?) AND
               	comments.created > ?
			GROUP BY comments.parent
			`
	rows, err := db.Query(sql, uri, mode, mode, after)
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

	comments := []Comment{}
	for rows.Next() {
		c := Comment{}

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
