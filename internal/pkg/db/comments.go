package db

import (
	"errors"
	"regexp"
	"time"

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

// Verify check comment invalid or valid
func (c *Comment) Verify() error {
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
		emailRegex := "[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?"
		ok, _ := regexp.MatchString(emailRegex, c.email.String)
		if !ok {
			return errors.New("invalid email")
		}
	}

	if c.Website.Valid {
		if len(c.email.String) > 254 {
			return errors.New("arbitrary length limit")
		}
		WebsiteRegex := "[(http(s)?):\\/\\/(www\\.)?a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)"

		ok, _ := regexp.MatchString(WebsiteRegex, c.Website.String)
		if !ok {
			return errors.New("invalid website address")
		}
	}

	return nil
}
