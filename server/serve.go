package server

import (
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"wrong.wang/x/go-isso/config"
	"wrong.wang/x/go-isso/database"
	"wrong.wang/x/go-isso/isso"
	"wrong.wang/x/go-isso/logger"
)

// Serve starts a new HTTP server.
func Serve(cfg config.Config) *http.Server {
	server := &http.Server{
		Handler:        setupHandler(cfg),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	switch {
	case strings.HasPrefix(cfg.Server.Listen, "unix://"):
		startUnixSocketServer(server, strings.TrimPrefix(cfg.Server.Listen, "unix://"))
	case strings.HasPrefix(cfg.Server.Listen, "http://"):
		server.Addr = strings.TrimPrefix(cfg.Server.Listen, "http://")
		startHTTPServer(server)
	default:
		logger.Fatal("not supported listen address:", cfg.Server.Listen)
	}
	return server
}

func startUnixSocketServer(server *http.Server, socketFile string) {
	os.Remove(socketFile)

	go func(sock string) {
		listener, err := net.Listen("unix", sock)
		if err != nil {
			logger.Fatal(`Server failed to start: %v`, err)
		}
		defer listener.Close()

		if err := os.Chmod(sock, 0666); err != nil {
			logger.Fatal(`Unable to change socket permission: %v`, err)
		}

		logger.Info(`Listening on Unix socket %q`, sock)
		if err := server.Serve(listener); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}(socketFile)
}

func startHTTPServer(server *http.Server) {
	go func() {
		logger.Info(`Listening on %q without TLS`, server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}()
}

func setupHandler(cfg config.Config) http.Handler {
	router := mux.NewRouter()
	router = router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		origin := isso.FindOrigin(r)
		for _, allowHost := range cfg.Host {
			if origin == allowHost {
				return true
			}
		}
		return false
	}).Subrouter()

	storage, err := database.New(cfg.DBPath, 1*time.Second)
	if err != nil {
		logger.Fatal("init database failed %w", err)
	}
	registerRoute(router, isso.New(cfg, storage))

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Host,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Origin", "Referer", "Content-Type"},
		ExposedHeaders:   []string{"X-Set-Cookie", "Date"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
		Debug:            false,
	})

	return c.Handler(router)
}
