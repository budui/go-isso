package httpd

import (
	"net/http"

	"github.com/gorilla/mux"
	"wrong.wang/x/go-isso/httpd/response"
	"wrong.wang/x/go-isso/isso"
)

func registerRoute(router *mux.Router, isso *isso.ISSO) {
	var h handler
	// single comment
	router.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		builder := response.New(w, r)
		isso.CreateComment(builder, r)
	}).Queries("uri", "{uri}").Methods("POST").Name("new")
	router.HandleFunc("/id/{id:[0-9]+}", h.WorkInProcess).Methods("GET").Name("view")
	router.HandleFunc("/id/{id:[0-9]+}", h.WorkInProcess).Methods("PUT").Name("edit")
	router.HandleFunc("/id/{id:[0-9]+}", h.WorkInProcess).Methods("DELETE").Name("delete")
	router.HandleFunc("/id/{id:[0-9]+}/like", h.WorkInProcess).Methods("POST").Name("like")
	router.HandleFunc("/id/{id:[0-9]+}/dislike", h.WorkInProcess).Methods("POST").Name("dislike")
	router.HandleFunc("/id/{id:[0-9]+}/{action:[edit|activate|delete]}/{key}", h.WorkInProcess).
		Methods("GET").Name("moderate_get")
	router.HandleFunc("/id/{id:[0-9]+}/{action:[edit|activate|delete]}>/{key}", h.WorkInProcess).
		Methods("POST").Name("moderate_post")
	router.HandleFunc("/id/{id:[0-9]+}/unsubscribe/{email}/{key}>", h.WorkInProcess).
		Methods("GET").Name("unsubscribe")

	// functional
	router.HandleFunc("/demo", h.WorkInProcess).Methods("GET").Name("demo")
	router.HandleFunc("/preview", h.WorkInProcess).Methods("POST").Name("preview")

	// amdin staff
	router.HandleFunc("/admin", h.WorkInProcess).Methods("GET").Name("admin")
	router.HandleFunc("/login", h.WorkInProcess).Methods("POST").Name("login")

	// ping
	router.HandleFunc("/ping", h.Ping).Name("ping")

	// total staff
	router.HandleFunc("/latest", h.WorkInProcess).Methods("GET").Name("latest")
	router.HandleFunc("/count", h.WorkInProcess).Methods("GET").Name("count")
	router.HandleFunc("/count", h.WorkInProcess).Methods("POST").Name("counts")
	router.HandleFunc("/", h.WorkInProcess).Methods("GET").Name("fetch")
}
