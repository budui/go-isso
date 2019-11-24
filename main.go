package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/budui/go-isso/config"
	"github.com/budui/go-isso/isso"
)

func main() {
	var (
		flagVersion    bool
		flagConfigFile string
	)
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("\tgo-isso [-v] -c <CONFIG PATH> [import|run] \n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&flagConfigFile, "c", "", "Load configuration file")
	flag.BoolVar(&flagVersion, "v", false, "Show application version")
	flag.Parse()

	if flagVersion {
		fmt.Printf("version        : %s\n", isso.Version)
		fmt.Printf("build timestamp: %s\n", isso.BuildTime)
		return
	}

	if flagConfigFile == "" {
		fmt.Printf("must specify configuration file\n\n")
		flag.Usage()
		return
	}

	cfg, err := config.Parse(flagConfigFile)
	if err != nil {
		fmt.Printf("[ERROR] read config file failed:\n\t%v\n", err)
		return
	}

	if flag.NArg() != 1 {
		fmt.Printf("need one and only one argument to spectify action.\n\n")
		flag.Usage()
		return
	}

	switch action := flag.Arg(0); action {
	case "import":
		fmt.Println("[ERROR] import still work in progress")
	case "run":
		issoWorker, err := isso.NewWorker(cfg)
		if err != nil {
			fmt.Printf("[ERROR] isso server initialization failed:\n\t%v\n", err)
			return
		}
		issoWorker.Run()
	default:
		fmt.Printf("%s is not supported action argument \n\n", action)
		flag.Usage()
	}
}
