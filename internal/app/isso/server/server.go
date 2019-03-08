package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/RayHY/go-isso/internal/app/isso/way"
	"github.com/RayHY/go-isso/internal/pkg/conf"
	"github.com/RayHY/go-isso/internal/pkg/db"
	log "github.com/RayHY/go-isso/internal/pkg/dlog"
)

// Server is the main struct for isso.
// it keep the all shared dependencies.
type Server struct {
	Router *way.Router
	db     db.Accessor
	Conf   conf.Config
	log    *log.Logger
}

// NewServer make a new Server
func NewServer(config conf.Config, inDebugMode bool) (*Server, error) {
	accessor, err := db.NewAccessor(config.Database, config.Guard)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Router: way.NewRouter(),
		Conf:   config,
		log:    log.New(os.Stdout, "", log.LstdFlags, inDebugMode),
		db:     accessor,
	}

	s.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, http.StatusText(404), 404)
	})

	err = s.registerRouters()

	if err != nil {
		return nil, fmt.Errorf("register routers failed - %v", err)
	}

	return s, nil
}

// Run the Server
func (s *Server) Run() error {
	s.log.Infof("Listening and serving HTTP on %s\n", s.Conf.Listen)
	// TODO middleware wrap s.Router.
	Router := s.LoggingHandler(s.Router)
	return http.ListenAndServe(strings.TrimPrefix(s.Conf.Listen, "http://"), Router)
}

// Close clear all res used by server.
func (s *Server) Close() {
	_ = s.db.Close()
}
