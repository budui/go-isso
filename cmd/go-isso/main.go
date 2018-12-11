package main

import (
	"flag"

	"log"

	"github.com/RayHY/go-isso/internal/app/isso/server"
	"github.com/RayHY/go-isso/internal/pkg/conf"
)

func main() {
	log.SetFlags(log.LstdFlags)
	configPath := flag.String("c", "./configs/go-isso.toml", "set configuration file")
	flag.Parse()

	config, err := conf.Load(*configPath)
	if err != nil {
		log.Fatalf("[FATA] Load Config Failed %v", err)
	}

	isso, err := server.NewServer(config)
	defer isso.Close()
	if err != nil {
		log.Fatalf("[FATA] Start Server Failed %v", err)
	}

	if err := isso.Run(); err != nil {
		log.Fatalf("[FATA] Could not start server: %s\n", err.Error())
	}
}
