package isso

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"wrong.wang/x/go-isso/logger"
	"wrong.wang/x/go-isso/response"
	"wrong.wang/x/go-isso/response/json"
	"wrong.wang/x/go-isso/validator"
)

// CreateComment create a new comment
func (isso *ISSO) CreateComment(rb response.Builder, req *http.Request) {
	commentWebsite := findOrigin(req)
	if commentWebsite == "" {
		json.BadRequest(rb, errors.New("can not find origin"))
		return
	}
	var comment submittedComment
	err := jsonBind(req.Body, &comment)
	if err != nil {
		json.BadRequest(rb, err)
		return
	}
	comment.URI = mux.Vars(req)["uri"]
	comment.RemoteAddr = findClientIP(req)

	if err := validator.Validate(comment); err != nil {
		json.BadRequest(rb, err)
		return
	}
	pretty.Println(comment)

	var thread Thread
	if isso.storage.ContainsThread(req.Context(), comment.URI) {
		thread, err = isso.storage.GetThreadByURI(req.Context(), comment.URI)
		if err != nil {
			json.ServerError(rb, fmt.Errorf("can not get thread %w", err))
			return
		}
	} else {
		thread, err = isso.storage.NewThread(req.Context(), comment.URI, comment.Title, commentWebsite)
		if err != nil {
			json.ServerError(rb, fmt.Errorf("can not create new thread %w", err))
			return
		}
	}

	if isso.config.Moderation.Enable {
		if isso.config.Moderation.ApproveAcquaintance &&
			comment.Email != nil &&
			isso.storage.IsApprovedAuthor(req.Context(), *comment.Email) {
			comment.Mode = 1
		} else {
			comment.Mode = 2
		}
	} else {
		comment.Mode = 1
	}

	c, err := isso.storage.NewComment(req.Context(), comment.Comment, thread.ID, comment.RemoteAddr)
	if err != nil {
		json.ServerError(rb, fmt.Errorf("can not create new comment %w", err))
		return
	}

	logger.Debug(fmt.Sprintf("new comment: %# v", pretty.Formatter(c)))

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

	if c.Mode == 2 {
		json.Accepted(rb, c)
	} else {
		json.Created(rb, c)
	}
}
