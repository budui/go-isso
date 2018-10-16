package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/jinxiapu/go-isso/internal/app/isso/server"
	"github.com/jinxiapu/go-isso/internal/pkg/conf"
)

func main() {
	configPath := flag.String("c", "../../configs/isso.conf", "set configuration file")
	flag.Parse()
	
	cfg, _ := conf.Load(*configPath)
	isso := server.NewServer(cfg)
	defer isso.Close()
	
	if err := isso.Run(); err != nil {
		logrus.Fatalf("Could not start server: %s\n", err.Error())
	}
}
