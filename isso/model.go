package isso

// Thread is comments thread
type Thread struct {
	ID    int64
	URI   string `validate:"required,uri"`
	Title string
}

// Comment is comment saved in database
type Comment struct {
	ID           int64    `json:"id"`
	Parent       *int64   `json:"parent"`
	Created      float64  `json:"created"`
	Modified     *float64 `json:"modified"`
	Mode         int      `json:"mode"`
	Text         string   `json:"text"  validate:"required,gte=3,lte=65535"`
	Author       string   `json:"author"  validate:"required,gte=1,lte=15"`
	Email        *string  `json:"email"  validate:"omitempty,email"`
	Website      *string  `json:"website"  validate:"omitempty,url"`
	Likes        int      `json:"like"`
	Dislikes     int      `json:"dislike"`
	Notification int      `json:"notification" validate:"omitempty,min=0,max=2"`
	RemoteAddr   string   `json:"-" validate:"required,ip"`
}

type submittedComment struct {
	Comment
	URI   string `json:"-" validate:"required,uri"`
	Title string `json:"title" validate:"omitempty"`
}

// Hash use provited hash worker to hash itself
func (c Comment) Hash(worker func(string) string) string {
	var s string
	if c.Email != nil {
		s = *c.Email
	} else {
		s = c.RemoteAddr
	}
	return worker(s)
}
