package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/RayHY/go-isso/internal/app/isso/service"
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
	hw     *service.HashWorker
	mdc    *service.MDConverter
}

// NewServer make a new Server
func NewServer(config conf.Config, inDebugMode bool) (*Server, error) {
	accessor, err := db.NewAccessor(config.Database.Sqlite.Path, config.Guard)
	if err != nil {
		return nil, err
	}
	hw, err := service.NewHashWorker(config.Hash.Algorithm, config.Hash.Salt)
	if err != nil {
		return nil, err
	}
	mdc := service.NewMDConverter(config.Markup)

	s := &Server{
		Router: way.NewRouter(),
		Conf:   config,
		log:    log.New(os.Stdout, "", log.LstdFlags, inDebugMode),
		db:     accessor,
		hw:     hw,
		mdc:    mdc,
	}

	s.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, http.StatusText(404), 404)
	})

	s.registerRouters()

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
