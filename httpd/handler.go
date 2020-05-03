package httpd

import (
	"net/http"

	"wrong.wang/x/go-isso/httpd/response"
	"wrong.wang/x/go-isso/isso"
)

type handler struct {
	isso isso.ISSO
}

func (h handler) WorkInProcess(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("work in process\n"))
}

func (h handler) newComment(w http.ResponseWriter, r *http.Request) {
	builder := response.New(w, r)
	h.isso.CreateComment(builder, r)
}

func (h handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}
