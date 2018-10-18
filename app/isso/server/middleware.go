package server

import (
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingHandler return a http.Handler that wraps h and logs requests
func (s *Server) LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)
		s.log.Debugf("%s - %s \"%s %s\" %d %s", r.RemoteAddr, start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}
