package main

import (
	"flag"

	"log"

	"github.com/RayHY/go-isso/internal/app/isso/server"
	"github.com/RayHY/go-isso/internal/pkg/conf"
)

func main() {
	inDebugMode := flag.Bool("d", false, "run go-isso in debug mode")
	configPath := flag.String("c", "./configs/go-isso.toml", "set configuration file")
	flag.Parse()

	if *inDebugMode {
		log.Print("[INFO] RUN IN DEBUG MODE.")
	}

	config, err := conf.Load(*configPath)
	if err != nil {
		log.Fatalf("[FATA] Load Config Failed : %v", err)
	}

	isso, err := server.NewServer(config, *inDebugMode)
	defer isso.Close()
	if err != nil {
		log.Fatalf("[FATA] Failed to Setup Application : %v", err)
	}

	if err := isso.Run(); err != nil {
		log.Fatalf("[FATA] Could Not Run Server: %s", err.Error())
	}
}
