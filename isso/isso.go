package isso

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/budui/go-isso/config"
	"github.com/budui/go-isso/pkg/logger"
)

// Worker is the main struct for go-isso,
// store all important struct like config, database,etc.
type Worker struct {
	logger *logger.Logger
	config *config.Config
	router *mux.Router
	server *http.Server
}

// NewWorker return a Worker. if any fields initiazed failed, return error.
func NewWorker(conf *config.Config) (*Worker, error) {
	log := logger.New(os.Stderr, "", logger.LstdFlags, conf.Debug)
	log.Debug("isso.Logger prepare ok.")

	r := mux.NewRouter()
	server := &http.Server{
		// TODO: remove this fallback string.
		Addr:         conf.Server.Listen,
		WriteTimeout: time.Second * 3,
		ReadTimeout:  time.Second * 3,
		IdleTimeout:  time.Second * 3,
		Handler:      r,
	}
	return &Worker{logger: log, config: conf, router: r, server: server}, nil
}

// Run start the daemon process for go-isso
func (ws *Worker) Run() {
	ws.logger.Debug("start run isso server")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	ws.logger.Printf("begin to listen on %v\n", ws.server.Addr)
	ws.server.ListenAndServe()

	<-stop
	ws.logger.Println("shutting down the process...")
}
