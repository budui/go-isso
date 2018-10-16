package server

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jinxiapu/go-isso/internal/app/isso/sender"
	"github.com/jinxiapu/go-isso/internal/pkg/conf"
	"github.com/jinxiapu/go-isso/internal/pkg/db"
	"github.com/jinxiapu/go-isso/pkg/way"
)

// Server is the main struct for isso.
// it keep the all shared dependencies.
type Server struct {
	Router *way.Router
	Sender *sender.Sender
	db     db.Worker
	Conf   conf.Configure
	log    *logrus.Logger
}

// NewServer make a new Server
func NewServer(config conf.Configure) *Server {
	s := &Server{
		Router: way.NewRouter(),
		Conf:   config,
		log:    logrus.New(),
		db:     db.NewWorker(config.Section("general").Key("dbpath").String(), db.Guard{}),
	}
	s.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, http.StatusText(404), 404)
	})
	s.registerRouters()
	//s.log.SetOutput()
	return s
}

// Run the Server
func (s *Server) Run() error {
	s.db.Prepare()
	listenAddr := s.Conf.Section("server").Key("listen").String()
	s.log.Infof("Listening and serving HTTP on %s\n", listenAddr)
	// TODO middleware wrap s.Router.
	return http.ListenAndServe(strings.TrimPrefix(listenAddr, "http://"), s.Router)
}

// Close clear all res used by server.
func (s *Server) Close() {
	
}
