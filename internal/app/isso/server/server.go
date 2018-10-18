package server

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jinxiapu/go-isso/internal/app/isso/sender"
	"github.com/jinxiapu/go-isso/internal/app/isso/util"
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
	hasher util.Hasher
}

// NewServer make a new Server
func NewServer(config conf.Configure, isDebug bool) *Server {
	hasher, err := util.NewHasher(config.Section("hash").Key("algorithm").String(), config.Section("hash").Key("salt").String())
	if err != nil {
		logrus.Fatal(err)
	}
	s := &Server{
		Router: way.NewRouter(),
		Conf:   config,
		log:    logrus.New(),
		db:     db.NewWorker(config.Section("general").Key("dbpath").String(), db.Guard{}),
		hasher: hasher,
	}
	s.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, http.StatusText(404), 404)
	})
	s.registerRouters()
	if isDebug {
		s.log.SetLevel(logrus.DebugLevel)
	}
	return s
}

// Run the Server
func (s *Server) Run() error {
	err := s.db.PrepareToWork()
	if err != nil {
		return err
	}
	listenAddr := s.Conf.Section("server").Key("listen").String()
	s.log.Infof("Listening and serving HTTP on %s\n", listenAddr)
	// TODO middleware wrap s.Router.
	Router := s.LoggingHandler(s.Router)
	return http.ListenAndServe(strings.TrimPrefix(listenAddr, "http://"), Router)
}

// Close clear all res used by server.
func (s *Server) Close() {

}
