package db

import (
	"encoding/json"

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

// To parse and unparse this JSON data, add this code to your project and do:
//
//    comment, err := unmarshalComment(bytes)
//    bytes, err = comment.Marshal()

func unmarshalComment(data []byte) (Comment, error) {
	var r Comment
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Comment) uarshal() ([]byte, error) {
	return json.Marshal(r)
}
