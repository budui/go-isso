package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/RayHY/go-isso/internal/app/isso/server"
	"github.com/RayHY/go-isso/internal/pkg/conf"
)

func main() {
	isDebug := flag.Bool("d", false, "run for debug")
	configPath := flag.String("c", "./configs/isso.conf", "set configuration file")
	flag.Parse()

	cfg, _ := conf.Load(*configPath)
	isso := server.NewServer(cfg, *isDebug)
	defer isso.Close()

	if err := isso.Run(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
