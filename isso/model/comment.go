package model

import "time"

// ReplyComment to reply
type ReplyComment struct {
	Dislike       int        `json:"dislike"`
	Like          int        `json:"like"`
	ID            int64        `json:"id"`
	Mode          int        `json:"mode"`
	Hash          string     `json:"hash"`
	Created       time.Time  `json:"created"`
	Modified      *time.Time `json:"modified"`
	GravatarImage string     `json:"gravatar_image"`

	AcceptComment
}

// AcceptComment contains fields that can be submitted
type AcceptComment struct {
	Text         string  `json:"text"  validate:"required,gte=3,lte=65535"`
	Author       string  `json:"author"  validate:"required,gte=1,lte=15"`
	Email        string  `json:"email"  validate:"required,email"`
	Website      *string `json:"website"  validate:"omitempty,url"`
	Parent       int     `json:"parent" validate:"omitempty"`
	Notification bool    `json:"notification" validate:"omitempty,min=0,max=2"`
	Title        string  `json:"title" validate:"omitempty"`
	
	URI          string  `json:"-" validate:"required,uri"`
	Mode         int     `json:"-"`
	RemoteAddr   string  `json:"-" validate:"required,ip"`
	ThreadID     int     `json:"-"`
}
