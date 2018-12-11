package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RayHY/go-isso/internal/app/isso/way"
	"github.com/RayHY/go-isso/internal/pkg/conf"
	"github.com/RayHY/go-isso/internal/pkg/db"
)

// Server is the main struct for isso.
// it keep the all shared dependencies.
type Server struct {
	Router *way.Router
	db     db.Worker
	Conf   conf.Config
	log    *log.Logger
}

// NewServer make a new Server
func NewServer(config conf.Config) *Server {
	s := &Server{
		Router: way.NewRouter(),
		Conf:   config,
		log:    log.New(os.Stdout, "", log.LstdFlags),
		db:     db.NewWorker(config.Database.Sqlite.Path, db.Guard{}),
	}
	s.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, http.StatusText(404), 404)
	})
	s.registerRouters()

	return s
}

// Run the Server
func (s *Server) Run() error {
	err := s.db.PrepareToWork()
	if err != nil {
		return err
	}
	listenAddr := s.Conf.Listen
	s.log.Printf("[INFO] Listening and serving HTTP on %s\n", listenAddr)
	// TODO middleware wrap s.Router.
	Router := s.LoggingHandler(s.Router)
	return http.ListenAndServe(strings.TrimPrefix(listenAddr, "http://"), Router)
}

// Close clear all res used by server.
func (s *Server) Close() {

}
