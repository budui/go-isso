package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/response"
)

func wrap(f func(response.Builder, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := newBuilder(w, r)
		f(b, r)
	}
}

func workInProcess(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("work in process\n"))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}

func registerRoute(router *mux.Router, isso *isso.ISSO) {
	// single comment
	router.HandleFunc("/new", wrap(isso.CreateComment)).Queries("uri", "{uri}").Methods("POST").Name("new")
	router.HandleFunc("/id/{id:[0-9]+}", workInProcess).Methods("GET").Name("view")
	router.HandleFunc("/id/{id:[0-9]+}", workInProcess).Methods("PUT").Name("edit")
	router.HandleFunc("/id/{id:[0-9]+}", workInProcess).Methods("DELETE").Name("delete")
	router.HandleFunc("/id/{id:[0-9]+}/like", workInProcess).Methods("POST").Name("like")
	router.HandleFunc("/id/{id:[0-9]+}/dislike", workInProcess).Methods("POST").Name("dislike")
	router.HandleFunc("/id/{id:[0-9]+}/{action:[edit|activate|delete]}/{key}", workInProcess).
		Methods("GET").Name("moderate_get")
	router.HandleFunc("/id/{id:[0-9]+}/{action:[edit|activate|delete]}>/{key}", workInProcess).
		Methods("POST").Name("moderate_post")
	router.HandleFunc("/id/{id:[0-9]+}/unsubscribe/{email}/{key}>", workInProcess).
		Methods("GET").Name("unsubscribe")

	// functional
	router.HandleFunc("/demo", workInProcess).Methods("GET").Name("demo")
	router.HandleFunc("/preview", workInProcess).Methods("POST").Name("preview")

	// amdin staff
	router.HandleFunc("/admin", workInProcess).Methods("GET").Name("admin")
	router.HandleFunc("/login", workInProcess).Methods("POST").Name("login")

	// ping
	router.HandleFunc("/ping", ping).Name("ping")

	// total staff
	router.HandleFunc("/latest", workInProcess).Methods("GET").Name("latest")
	router.HandleFunc("/count", workInProcess).Methods("GET").Name("count")
	router.HandleFunc("/count", wrap(isso.CountComment())).Methods("POST").Name("counts")
	router.HandleFunc("/js/embed.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		http.ServeFile(w, r, "./static/js/embed.min.js")
	})

	router.HandleFunc("/", wrap(isso.FetchComments())).Queries("uri", "{uri}").Methods("GET").Name("fetch")
}
