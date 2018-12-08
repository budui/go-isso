package server

func (s *Server) registerRouters() {
	s.Router.NotFound = s.handleStatusCode(404)
	s.Router.HandleFunc("GET", "/hello/:name", s.handleHello())
	s.Router.HandleFunc("GET", "/", s.handleFetch())
	s.Router.HandleFunc("POST", "/new", s.handleNew())
	s.Router.HandleFunc("GET", "/count", s.handleStatusCode(501))
	s.Router.HandleFunc("POST", "/count", s.handleStatusCode(501))
	s.Router.HandleFunc("GET", "/feed", s.handleStatusCode(501))
	s.Router.HandleFunc("GET", "/id/:id", s.handleStatusCode(501))
	s.Router.HandleFunc("PUT", "/id/:id", s.handleStatusCode(501))
	s.Router.HandleFunc("DELETE", "/id/:id", s.handleStatusCode(501))
}
