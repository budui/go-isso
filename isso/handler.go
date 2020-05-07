package isso

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"wrong.wang/x/go-isso/isso/model"
	"wrong.wang/x/go-isso/isso/request"
	"wrong.wang/x/go-isso/logger"
	"wrong.wang/x/go-isso/response"
	"wrong.wang/x/go-isso/response/json"
)

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
		model.SubmitComment
		Created string
	}{comment, strconv.FormatInt(time.Now().Unix(), 10)})
}
