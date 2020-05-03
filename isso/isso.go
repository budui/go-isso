package isso

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/kr/pretty"
	"wrong.wang/x/go-isso/config"
	"wrong.wang/x/go-isso/isso/model"
	"wrong.wang/x/go-isso/isso/request"
	"wrong.wang/x/go-isso/isso/response"
	"wrong.wang/x/go-isso/isso/response/json"
	"wrong.wang/x/go-isso/logger"
)

// ISSO do the main logical staff
type ISSO struct {
	storage Storage
	config  config.Config
	guard   guard
}

type guard struct {
	v  *Validator
	sc *securecookie.SecureCookie
}

// CreateComment create a new comment
func (isso *ISSO) CreateComment(rb response.Builder, req *http.Request) {
	comment, err := decodeAcceptComment(req.Body)
	if err != nil {
		json.BadRequest(rb, err)
		return
	}
	comment.URI = mux.Vars(req)["uri"]
	comment.RemoteAddr = request.FindClientIP(req)

	if err := isso.guard.v.Validate(comment); err != nil {
		json.BadRequest(rb, err)
		return
	}

	var thread model.Thread
	if isso.storage.ContainsThread(comment.URI) {
		thread, err = isso.storage.GetThreadByURI(comment.URI)
		if err != nil {
			json.ServerError(rb, fmt.Errorf("can not get thread %w", err))
			return
		}
	} else {
		thread, err = isso.storage.NewThread(comment.URI, comment.Title, req.Host)
		if err != nil {
			json.ServerError(rb, fmt.Errorf("can not create new thread %w", err))
			return
		}
	}

	comment.ThreadID = thread.ID
	if isso.config.Moderation.Enable {
		if isso.config.Moderation.ApproveAcquaintance && isso.storage.IsApprovedAuthor(comment.Email) {
			comment.Mode = 1
		} else {
			comment.Mode = 2
		}
	} else {
		comment.Mode = 1
	}

	c, err := isso.storage.NewComment(comment)
	if err != nil {
		json.ServerError(rb, fmt.Errorf("can not create new comment %w", err))
		return
	}

	logger.Debug(fmt.Sprintf("new comment: %# v", pretty.Formatter(comment)))

	if encoded, err := isso.guard.sc.Encode(fmt.Sprintf("%v", c.ID),
		map[int64][20]byte{c.ID: sha1.Sum([]byte(c.Text))}); err == nil {
		cookie := &http.Cookie{
			Name:   fmt.Sprintf("%v", c.ID),
			Value:  encoded,
			Path:   "/",
			MaxAge: isso.config.MaxAge,
		}
		if v := cookie.String(); v != "" {
			rb.WithHeader("Set-Cookie", v)
		}
	}

	json.Created(rb, struct {
		model.AcceptComment
		Created string
	}{comment, strconv.FormatInt(time.Now().Unix(), 10)})
}

// New a ISSO instance
func New(cfg config.Config, storage Storage) *ISSO {
	HashKey, err := storage.GetPreference("hask-key")
	if err != nil {
		HashKey = string(securecookie.GenerateRandomKey(64))
	}
	BlockKey, err := storage.GetPreference("block-key")
	if err != nil {
		BlockKey = string(securecookie.GenerateRandomKey(32))
	}
	return &ISSO{
		config: cfg,
		guard: guard{
			v:  NewValidator(),
			sc: securecookie.New([]byte(HashKey), []byte(BlockKey)),
		},
		storage: storage,
	}
}
