package cli

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"wrong.wang/x/go-isso/config"
	"wrong.wang/x/go-isso/httpd"
	"wrong.wang/x/go-isso/logger"
)

func startDaemon(cfg config.Config) {
	logger.Info("Starting go-isso...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	go showProcessStatistics()

	var httpServer *http.Server

	httpServer = httpd.Serve(cfg)

	<-stop
	logger.Info("Shutting down the process...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if httpServer != nil {
		httpServer.Shutdown(ctx)
	}

	logger.Info("Process gracefully stopped")
}

func showProcessStatistics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		logger.Debug("Sys=%vK, InUse=%vK, HeapInUse=%vK, StackSys=%vK, StackInUse=%vK, GoRoutines=%d, NumCPU=%d",
			m.Sys/1024, (m.Sys-m.HeapReleased)/1024, m.HeapInuse/1024, m.StackSys/1024, m.StackInuse/1024,
			runtime.NumGoroutine(), runtime.NumCPU())
		time.Sleep(30 * time.Second)
	}
}
