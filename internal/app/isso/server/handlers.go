package server

import (
	"encoding/json"
	"fmt"
	"github.com/RayHY/go-isso/internal/app/isso/service"
	"github.com/RayHY/go-isso/internal/pkg/dlog"
	"github.com/microcosm-cc/bluemonday"
	"html"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/RayHY/go-isso/internal/app/isso/way"
	"github.com/RayHY/go-isso/internal/pkg/db"
	"gopkg.in/guregu/null.v3"
)

func parseURLParam2type(Q url.Values, key string, item interface{}) error {
	r, ok := Q[key]
	switch value := item.(type) {
	case *null.String:
		if !ok {
			*value = null.NewString("", false)
			return nil
		}
		*value = null.StringFrom(r[0])
	case *null.Int:
		if !ok {
			*value = null.NewInt(0, false)
			return nil
		}
		v, err := strconv.Atoi(r[0])
		if err != nil || v < 0 {
			*value = null.NewInt(0, false)
			return fmt.Errorf("param '%s' invalid", key)
		}
		*value = null.IntFrom(int64(v))
	case *float64:
		if !ok {
			*value = 0.00
			return nil
		}
		v, err := strconv.ParseFloat(r[0], 64)
		if err != nil {
			*value = 0.00
			return fmt.Errorf("param '%s' invalid", key)
		}
		*value = v
	default:
		panic(fmt.Sprintf("do not support type : %T", value))
	}
	return nil
}

func sanitizeUserInput(in null.String) null.String {
	if !in.Valid {
		return in
	}
	return null.StringFrom(html.EscapeString(bluemonday.UGCPolicy().Sanitize(in.String)))
}

func getAPIExceptionHandler(logger *dlog.Logger, APIName string) func(http.ResponseWriter, string, error, int) {
	return func(w http.ResponseWriter, info string, err error, code int) {
		levelName := "ERROR"
		if code < http.StatusInternalServerError {
			levelName = "WARN"
			jsonError(w, info, code)
		} else {
			http.Error(w, http.StatusText(code), code)
		}

		if err != nil {
			logger.Printf("[%s] @api.%s: %s - %v", levelName, APIName, info, err)
		} else {
			logger.Printf("[%s] @api.%s: %s", levelName, APIName, info)
		}
	}
}

func jSON(w http.ResponseWriter, v interface{}, status int) {
	// Set before WriteHeader cause https://golang.org/pkg/net/http/?#ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(v)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jSON(w, map[string]string{"error": message}, status)
}

func (s *Server) handleStatusCode(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, http.StatusText(code), code)
	}
}

// example 'https://comments.example.com/?uri=/thread/&limit=2&nested_limit=5'
func (s *Server) handleFetch() http.HandlerFunc {
	ExceptionHandler := getAPIExceptionHandler(s.log, "Fetch")
	type reply struct {
		db.Comment
		HiddenReplies *int64  `json:"hidden_replies,omitempty"`
		TotalReplies  *int64  `json:"total_replies,omitempty"`
		Replies       []reply `json:"replies"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var uri null.String
		var limit, parent, nestedLimit, plain null.Int
		var after float64

		URLParamsPair := map[string]interface{}{
			"uri": &uri, "limit": &limit, "parent": &parent,
			"nested_limit": &nestedLimit, "plain": &plain, "after": &after,
		}

		for k, v := range URLParamsPair {
			err := parseURLParam2type(r.URL.Query(), k, v)
			if err != nil {
				ExceptionHandler(w, fmt.Sprintf("parse param '%s' failed", k), err, http.StatusBadRequest)
				return
			}
		}
		if !uri.Valid {
			ExceptionHandler(w, "missing uri query", nil, http.StatusBadRequest)
			return
		}

		if !(plain.Int64 == 0 || plain.Int64 == 1) {
			ExceptionHandler(w, "param 'plain' invalid", nil, http.StatusBadRequest)
			return
		}

		replyCounts, err := s.db.CountReply(uri.String, db.ModePublic, after)
		if err != nil {
			ExceptionHandler(w, "reply count failed", err, http.StatusInternalServerError)
			return
		}

		_, ok := replyCounts[parent]
		if !ok {
			replyCounts[parent] = 0
		}

		var rJSON struct {
			TotalReplies  int64    `json:"total_replies"`
			Replies       []reply  `json:"replies"`
			ID            null.Int `json:"id"`
			HiddenReplies int64    `json:"hidden_replies"`
		}
		rJSON.ID = parent
		rJSON.TotalReplies = replyCounts[parent]

		FetchComments := func(w http.ResponseWriter, parent, limit null.Int) ([]reply, error) {
			comments, err := s.db.Fetch(uri.String, db.ModePublic, after, parent, "id", true, limit)
			if err != nil {
				return nil, err
			}
			replies := []reply{}
			for _, c := range comments {
				c.Hash = s.hw.Hash(c.EmailOrIP())
				if !plain.Valid || plain.Int64 != 1 {
					c.Text = s.mdc.Run(c.Text)
				}
				r := reply{c, nil, nil, []reply{}}
				replies = append(replies, r)
			}
			return replies, nil
		}

		rJSON.Replies, err = FetchComments(w, parent, limit)
		if err != nil {
			ExceptionHandler(w, "fetch comment failed", err, http.StatusInternalServerError)
			return
		}

		// run only when parent == NULL
		// I don't understand why but isso run like this.
		// so just keep compatible with isso api.
		if !parent.Valid {
			for i := range rJSON.Replies {
				rJSON.Replies[i].TotalReplies = new(int64)
				rJSON.Replies[i].HiddenReplies = new(int64)
				count, ok := replyCounts[null.IntFrom(rJSON.Replies[i].ID)]
				if !ok {
					*rJSON.Replies[i].TotalReplies = 0
				} else {
					*rJSON.Replies[i].TotalReplies = count
					rJSON.Replies[i].Replies, err = FetchComments(w, null.IntFrom(rJSON.Replies[i].ID), nestedLimit)
					if err != nil {
						ExceptionHandler(w, "fetch comment failed", err, http.StatusInternalServerError)
						return
					}
				}
				*rJSON.Replies[i].HiddenReplies = *rJSON.Replies[i].TotalReplies - int64(len(rJSON.Replies[i].Replies))
			}
		}

		rJSON.HiddenReplies = replyCounts[parent] - int64(len(rJSON.Replies))
		jSON(w, rJSON, http.StatusOK)
	}
}

func (s *Server) handleNew() http.HandlerFunc {
	ExceptionHandler := getAPIExceptionHandler(s.log, "New")
	var mode int64
	var successCode int
	if s.Conf.Moderation.Enable {
		mode = 1
		successCode = 201
	} else {
		mode = 2
		successCode = 202
	}

	var titleExtractor = service.NewTitleExtractor(http.Client{Timeout: time.Second * 5})

	return func(w http.ResponseWriter, r *http.Request) {
		var uri null.String
		_ = parseURLParam2type(r.URL.Query(), "uri", &uri)
		if !uri.Valid {
			ExceptionHandler(w, "missing uri query", nil, http.StatusBadRequest)
			return
		}

		var nc struct {
			Text         string      `json:"text"`
			Parent       null.Int    `json:"parent"`
			Author       null.String `json:"author"`
			Email        null.String `json:"email"`
			Website      null.String `json:"website"`
			Title        null.String `json:"title"`
			Notification int64       `json:"notification"`
		}
		err := json.NewDecoder(http.MaxBytesReader(w, r.Body, int64(1<<14))).Decode(&nc)
		if err != nil {
			ExceptionHandler(w, "decode input json failed", err, http.StatusBadRequest)
			return
		}
		nc.Website = sanitizeUserInput(nc.Website)
		nc.Email = sanitizeUserInput(nc.Email)
		nc.Author = sanitizeUserInput(nc.Author)

		var thread db.Thread
		if ok, err := s.db.Contains(uri.String); err != nil {
			if ok {
				thread, _ = s.db.GetThreadWithURI(uri.String)
			} else {
				if !nc.Title.Valid {
					threadURL := path.Join(s.Conf.Hosts[0], uri.String)
					title, err := titleExtractor.Get(threadURL)
					if err != nil {
						ExceptionHandler(w, fmt.Sprintf("get thread page(%v) failed", threadURL),
							err, http.StatusNotFound)
						return
					}
					nc.Title.SetValid(title)
				}
				thread, err = s.db.NewThread(uri.String, nc.Title)
				if err != nil {
					ExceptionHandler(w, fmt.Sprintf("new thread(%v, %v) failed", uri.String, nc.Title),
						err, http.StatusInternalServerError)
					return
				}
			}
		} else {
			ExceptionHandler(w, fmt.Sprintf("check new uri(%s) failed", uri.String),
				err, http.StatusInternalServerError)
			return
		}

		c := db.NewComment(thread.ID, nc.Parent, mode, strings.Split(r.RemoteAddr, ":")[0], nc.Text,
			nc.Author, nc.Email, nc.Website, nc.Notification)
		if err := c.Verify(); err != nil {
			ExceptionHandler(w, "verify user input failed", err, http.StatusBadRequest)
			return
		}
		c, err = s.db.Add(uri.String, c)
		if err != nil {
			ExceptionHandler(w, "can not add comment into database", err, http.StatusInternalServerError)
			return
		}

		// TODO: session and cookies for editing comments.

		c.Hash = s.hw.Hash(c.EmailOrIP())
		c.Text = s.mdc.Run(c.Text)
		jSON(w, c, successCode)
	}
}

// example 'https://comments.example.com/id/4'
func (s *Server) handleView() http.HandlerFunc {
	ExceptionHandler := getAPIExceptionHandler(s.log, "View")
	return func(w http.ResponseWriter, r *http.Request) {
		CommentID, err := strconv.Atoi(way.Param(r.Context(), "id"))
		if err != nil {
			ExceptionHandler(w, "Invalid ID", err, http.StatusBadRequest)
			return
		}
		c, err := s.db.Get(int64(CommentID))
		if err != nil {
			ExceptionHandler(w, "Can't get corresponding comment", err, http.StatusBadRequest)
			return
		}

		var plain null.Int
		err = parseURLParam2type(r.URL.Query(), "plain", &plain)
		if err != nil {
			ExceptionHandler(w, "param 'plain' invalid", err, http.StatusBadRequest)
			return
		}
		if !plain.Valid || plain.Int64 != 1 {
			c.Text = s.mdc.Run(c.Text)
		}
		jSON(w, c, 200)
	}
}

// curl -X PUT 'https://comments.example.com/id/23' -d \
// {"text": "I see your point. However, I still disagree.", "website":\
// "maxrant.important.com"} -H 'Content-Type: application/json' -b cookie.txt
func (s *Server) handleEdit() http.HandlerFunc {
	ExceptionHandler := getAPIExceptionHandler(s.log, "Edit")
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: first check whether user can edit this comment.
		var nc struct {
			Text    string      `json:"text"`
			Author  null.String `json:"author"`
			Website null.String `json:"website"`
		}
		err := json.NewDecoder(http.MaxBytesReader(w, r.Body, int64(1<<14))).Decode(&nc)
		if err != nil {
			ExceptionHandler(w, "decode input json failed", err, http.StatusBadRequest)
			return
		}
		nc.Website = sanitizeUserInput(nc.Website)
		nc.Author = sanitizeUserInput(nc.Author)
		if len(nc.Text) < 3 || len(nc.Text) > 65535 {
			ExceptionHandler(w, "input text size too long or too small", nil, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(way.Param(r.Context(), "id"))
		c, err := s.db.Update(int64(id), nc.Text, nc.Author, nc.Website)
		if err != nil {
			ExceptionHandler(w, "can't update corresponding comment", err, http.StatusInternalServerError)
			return
		}

		var plain null.Int
		err = parseURLParam2type(r.URL.Query(), "plain", &plain)
		if err != nil {
			ExceptionHandler(w, "param 'plain' invalid", err, http.StatusBadRequest)
			return
		}
		if !plain.Valid || plain.Int64 != 1 {
			c.Text = s.mdc.Run(c.Text)
		}
		jSON(w, c, 200)
	}
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := way.Param(r.Context(), "name")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, fmt.Sprintf("hello, %s", name))
	}
}
