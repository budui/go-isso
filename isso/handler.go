package isso

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/kr/pretty"
	"wrong.wang/x/go-isso/logger"
	"wrong.wang/x/go-isso/response"
	"wrong.wang/x/go-isso/response/json"
	"wrong.wang/x/go-isso/validator"
)

// CreateComment create a new comment
func (isso *ISSO) CreateComment(rb response.Builder, req *http.Request) {
	commentWebsite := FindOrigin(req)
	if commentWebsite == "://" {
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
	thread, err = isso.storage.GetThreadByURI(req.Context(), comment.URI)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no thread realted to this uri
			// so create new thread
			if thread, err = isso.storage.NewThread(req.Context(), comment.URI, comment.Title, commentWebsite); err != nil {
				json.ServerError(rb, fmt.Errorf("can not create new thread %w", err))
				return
			}
		} else {
			// can not handled error
			json.ServerError(rb, fmt.Errorf("can not get thread %w", err))
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

// FetchComments fetch all related comments
func (isso *ISSO) FetchComments() func(rb response.Builder, req *http.Request) {
	type urlParm struct {
		Parent      *int64  `schema:"parent"`
		Limit       int64   `schema:"limit"`
		NestedLimit int64   `schema:"nested_limit"`
		After       float64 `schema:"after"`
		Plain       int64   `schema:"plain"`
	}
	type reply struct {
		Comment
		HiddenReplies *int64  `json:"hidden_replies,omitempty"`
		TotalReplies  *int64  `json:"total_replies,omitempty"`
		Replies       []reply `json:"replies"`
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	makeReplies := func(cs []Comment, after float64, limit int64, plain bool) []reply {
		var replies []reply
		var count int64
		if limit <= 0 {
			limit = int64(len(cs) + 1)
		}
		for _, c := range cs {
			if c.Created > after && count < limit {
				count++
				replies = append(replies, reply{c, nil, nil, nil})
			}
		}
		return replies
	}

	return func(rb response.Builder, req *http.Request) {
		var urlparm urlParm
		err := decoder.Decode(&urlparm, req.URL.Query())
		if err != nil {
			json.BadRequest(rb, err)
			return
		}
		var parent int64
		if urlparm.Parent == nil {
			parent = -1
		} else {
			parent = *urlparm.Parent
		}
		var plain bool
		if urlparm.Plain != 0 {
			plain = true
		}

		replyCount, err := isso.storage.CountReply(req.Context(), mux.Vars(req)["uri"], 5, urlparm.After)
		if err != nil {
			json.ServerError(rb, err)
			return
		}
		// param `after` may cause the loss of old comment's parent
		if _, ok := replyCount[parent]; !ok {
			replyCount[parent] = 0
		}

		commentsByParent, err := isso.storage.FetchCommentsByURI(req.Context(), mux.Vars(req)["uri"], parent, 5, "id", true)
		if err != nil {
			json.ServerError(rb, fmt.Errorf("fetch comments failed %w", err))
			return
		}
		rJSON := struct {
			TotalReplies  int64   `json:"total_replies"`
			Replies       []reply `json:"replies"`
			ID            *int64  `json:"id"`
			HiddenReplies int64   `json:"hidden_replies"`
		}{
			ID: urlparm.Parent,
		}

		// null parent, only fetch top-comment
		if parent == -1 {
			// parent == -1 means need all comment's, here TotalReplies means top-leval comments
			rJSON.TotalReplies = replyCount[0]

			rJSON.Replies = makeReplies(commentsByParent[0], urlparm.After, urlparm.Limit, plain)
			rJSON.HiddenReplies = rJSON.TotalReplies - int64(len(rJSON.Replies))
			var zero int64
			for i := range rJSON.Replies {
				count, ok := replyCount[rJSON.Replies[i].ID]
				if !ok {
					rJSON.Replies[i].TotalReplies = &zero
				} else {
					rJSON.Replies[i].TotalReplies = &count
					rJSON.Replies[i].Replies = makeReplies(commentsByParent[rJSON.Replies[i].ID], urlparm.After, urlparm.NestedLimit, plain)
					cc := *rJSON.Replies[i].TotalReplies - int64(len(rJSON.Replies[i].Replies))
					rJSON.Replies[i].HiddenReplies = &cc
				}
			}

		} else if parent > 0 {
			rJSON.TotalReplies = replyCount[parent]
			rJSON.Replies = makeReplies(commentsByParent[parent], urlparm.After, urlparm.Limit, plain)
			rJSON.HiddenReplies = rJSON.TotalReplies - int64(len(rJSON.Replies))
		} else {
			// parent = 0 not exist
			rJSON.TotalReplies = 0
			rJSON.Replies = []reply{}
			rJSON.HiddenReplies = 0
		}
		json.OK(rb, rJSON)
	}
}
