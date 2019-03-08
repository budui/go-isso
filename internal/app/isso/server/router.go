package server

import (
	"errors"
	"github.com/RayHY/go-isso/internal/app/isso/service"
	"github.com/gorilla/securecookie"
)

// registerRouters register routers for isso.
// also provide some services for the corresponding routes
func (s *Server) registerRouters() error {

	HashService, err := service.NewHashWorker(s.Conf.Hash.Algorithm, s.Conf.Hash.Salt)
	if err != nil {
		return errors.New("initial hash worker failed")
	}
	MDConverterService := service.NewMDConverter(s.Conf.Markup)

	cookieKey, err := s.db.GetPreference("session-key")
	if err != nil {
		return errors.New("get session-key failed")
	}
	secureCookieService := securecookie.New([]byte(cookieKey.String), nil)

	secureCookieService.MaxAge(int(s.Conf.Guard.EditMaxAge))

	s.Router.NotFound = s.handleStatusCode(404)
	s.Router.HandleFunc("GET", "/hello/:name", s.handleHello())
	s.Router.HandleFunc("GET", "/", s.handleFetch(MDConverterService, HashService))
	s.Router.HandleFunc("POST", "/new", s.handleNew(MDConverterService, HashService, secureCookieService))
	s.Router.HandleFunc("GET", "/count", s.handleStatusCode(501))
	s.Router.HandleFunc("POST", "/count", s.handleStatusCode(501))
	s.Router.HandleFunc("GET", "/feed", s.handleStatusCode(501))
	s.Router.HandleFunc("GET", "/id/:id", s.handleView(MDConverterService, HashService))
	s.Router.HandleFunc("PUT", "/id/:id", s.handleEdit(MDConverterService, HashService, secureCookieService))
	s.Router.HandleFunc("DELETE", "/id/:id", s.handleDelete())
	return nil
}
