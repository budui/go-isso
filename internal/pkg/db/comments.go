package db

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/RayHY/go-isso/internal/app/isso/service"
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
func NewComment(tid int64, Parent null.Int, mode int64, remoteAddr string,
	text string, author, email, website null.String, notification int64) Comment {
	c := Comment{
		tid:          tid,
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

	WebsiteRegex := `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`
	if len(c.Text) < 3 {
		return errors.New("text is too short (minimum length: 3)")
	}
	if len(c.Text) > 65535 {
		return errors.New("text is too long (maximum length: 65535)")
	}
	if c.Parent.Valid && c.Parent.Int64 <= 0 {
		return errors.New("parent must be an integer > 0")
	}

	if c.Author.Valid {
		c.Author.String = strings.TrimSpace(c.Author.String)
		if len(c.Author.String) > 63 {
			return errors.New("too long author name")
		}
	}

	if c.email.Valid {
		c.email.String = strings.TrimSpace(c.email.String)
		if len(c.email.String) > 254 {
			return errors.New("too long email")
		}
		ok, _ := regexp.MatchString(emailRegex, c.email.String)
		if !ok {
			return errors.New("invalid email")
		}
	}

	if c.Website.Valid {
		c.Website.String = strings.TrimSpace(c.Website.String)
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

// CookieValue return the value that can be used as edit cookie value.
func (c *Comment) CookieValue() map[int64][20]byte {
	return map[int64][20]byte{
		c.ID: sha1.Sum([]byte(c.Text)),
	}
}


func (db *database) Delete(id int64) (Comment, error) {
	return Comment{}, nil
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
		return nil, fmt.Errorf("db.CountReply: %v", err)
	}
	defer rows.Close()
	countResult := make(map[null.Int]int64)
	var parent null.Int
	var count int64
	for rows.Next() {
		err := rows.Scan(&parent, &count)
		if err != nil {
			return nil, fmt.Errorf("db.CountReply: %v", err)
		}
		countResult[parent] = count
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("db.CountReply: %v", err)
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
		return nil, fmt.Errorf("db.Fetch: %v", err)
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
			return nil, fmt.Errorf("db.Fetch: %v", err)
		}
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("db.Fetch: %v", err)
	}
	return comments, nil
}

// Add new comment to DB and return a complete Comment
func (db *database) Add(uri string, c Comment) (Comment, error) {
	if c.Parent.Valid {
		parent, err := db.Get(c.Parent.Int64)
		if err != nil {
			if err == sql.ErrNoRows {
				return Comment{}, newError("can't find specify parent")
			}
			return Comment{}, fmt.Errorf("db.Add: get parent failed. - %v", err)
		}
		parentThread, err := db.GetThreadWithID(parent.tid)
		if err != nil {
			return Comment{}, fmt.Errorf("db.Add: get parent's thread failed. - %v", err)
		}
		if parentThread.URI != uri {
			return Comment{}, newError("parent's thread and comment's thread are different")
		}
		if parent.Parent.Valid {
			c.Parent = parent.Parent
		}
	}

	c.voters = service.GenBloomfilterfunc(c.remoteAddr)

	stmt, err := db.Prepare(`
	INSERT INTO comments (
        	tid, parent, created, modified, mode, remote_addr,
			text, author, email, website, voters, notification
		)
    SELECT threads.id, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			FROM threads WHERE threads.uri = ?;`)
	if err != nil {
		return Comment{}, fmt.Errorf("db.Add: %v", err)
	}
	_, err = stmt.Exec(c.Parent, c.Created, c.Modified, c.Mode, c.remoteAddr,
		c.Text, c.Author, c.email, c.Website, c.voters, c.notification, uri)
	if err != nil {
		return Comment{}, fmt.Errorf("db.Add: %v", err)
	}

	err = db.QueryRow(`SELECT c.* FROM comments AS c INNER JOIN threads ON threads.uri = ? ORDER BY c.id DESC LIMIT 1`, uri).Scan(
		&c.tid, &c.ID, &c.Parent, &c.Created, &c.Modified,
		&c.Mode, &c.remoteAddr, &c.Text, &c.Author, &c.email, &c.Website,
		&c.Likes, &c.Dislikes, &c.voters, &c.notification,
	)
	if err != nil {
		return Comment{}, fmt.Errorf("db.Add: %v", err)
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
		if err == sql.ErrNoRows {
			return Comment{}, sql.ErrNoRows
		}
		return Comment{}, fmt.Errorf("db.Get: %v", err)
	}
	return c, nil
}

func (db *database) Update(id int64, text string, author, website null.String) (Comment, error) {
	stmt, err := db.Prepare(`UPDATE comments SET modified = ?, text = ?, author = ?, website = ? WHERE id=?;`)
	if err != nil {
		return Comment{}, fmt.Errorf("db.Update: %v", err)
	}
	_, err = stmt.Exec(float64(time.Now().UnixNano())/float64(1e9), text, author, website, id)
	if err != nil {
		return Comment{}, fmt.Errorf("db.Update: %v", err)
	}

	c, err := db.Get(id)
	if err != nil {
		return Comment{}, fmt.Errorf("db.Update: %v", err)
	}
	return c, nil
}
