package db

import (
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/RayHY/go-isso/internal/app/isso/util"
	"gopkg.in/guregu/null.v3"
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

// Comment represent replies
type Comment struct {
	tid          int64
	ID           int64      `json:"id"`
	Parent       null.Int   `json:"parent"`
	Created      float64    `json:"created"`
	Modified     null.Float `json:"modified"`
	Mode         int64      `json:"mode"`
	remoteAddr   string
	Text         string      `json:"text"`
	Author       null.String `json:"author"`
	email        null.String
	Website      null.String `json:"website"`
	Likes        int64       `json:"likes"`
	Dislikes     int64       `json:"dislikes"`
	voters       []byte
	notification int64
	// not store in database.
	Hash string `json:"hash"`
}

// NewComment return a new comment struct with setable fileds.
// ID, tid, or others generate later
func NewComment(Parent null.Int, mode int64, remoteAddr string,
	text string, author, email, website null.String, notification int64) Comment {
	c := Comment{
		Parent:       Parent,
		Mode:         mode,
		remoteAddr:   remoteAddr,
		Text:         text,
		Author:       author,
		email:        email,
		Website:      website,
		notification: notification,
	}
	c.Created = float64(time.Now().UnixNano()) / float64(1e9)
	return c
}

func (c *Comment) EmailOrIP() string {
	if c.email.Valid {
		return c.email.String
	}
	return c.remoteAddr
}

// Verify check comment invalid or valid
func (c *Comment) Verify() error {
	emailRegex := "[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?"
	WebsiteRegex := `[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	if len(c.Text) < 3 {
		return errors.New("text is too short (minimum length: 3)")
	}
	if len(c.Text) > 65535 {
		return errors.New("text is too long (maximum length: 65535)")
	}
	if c.Parent.Valid && c.Parent.Int64 <= 0 {
		return errors.New("parent must be an integer > 0")
	}

	if c.email.Valid {
		if len(c.email.String) > 254 {
			return errors.New("too long email")
		}
		ok, _ := regexp.MatchString(emailRegex, c.email.String)
		if !ok {
			return errors.New("invalid email")
		}
	}

	if c.Website.Valid {
		if len(c.email.String) > 254 {
			return errors.New("arbitrary length limit")
		}

		ok, _ := regexp.MatchString(WebsiteRegex, c.Website.String)
		if !ok {
			return errors.New("invalid website address")
		}
	}

	return nil
}

// commentAccessor defines all usual access ops avail for comment.
type commentAccessor interface {
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
// orderBy can be select from ['id', 'created', 'modified', 'likes', 'dislikes']
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

// Add new comment to DB and return a complete Comment
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
