package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"wrong.wang/x/go-isso/config"
	"wrong.wang/x/go-isso/logger"
	"wrong.wang/x/go-isso/version"
)

// Parse parses command line arguments
func Parse() {
	var (
		flagVersion        bool
		flagDebug          bool
		flagConfigFilePath string
	)
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("\tgo-isso [-v] -c <CONFIG PATH> [import|run] \n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&flagConfigFilePath, "c", "", "Load configuration file")
	flag.BoolVar(&flagVersion, "v", false, "Show application version")
	flag.BoolVar(&flagDebug, "d", false, "turn on debug mode")
	flag.Parse()

	if flagVersion {
		fmt.Println("Version:", version.Version)
		fmt.Println("Build Date:", version.BuildTime)
		fmt.Println("Go Version:", runtime.Version())
		fmt.Println("Compiler:", runtime.Compiler)
		fmt.Println("Arch:", runtime.GOARCH)
		fmt.Println("OS:", runtime.GOOS)
		return
	}

	if flagConfigFilePath == "" {
		fmt.Printf("must specify configuration file\n\n")
		flag.Usage()
		return
	}

	if flagDebug {
		logger.EnableDebug()
	}

	cfg, err := config.Parse(flagConfigFilePath)
	if err != nil {
		logger.Fatal("can not read config file: %v", err)
	}
	if cfg.LogFilePath != "" {
		logFile, err := os.Create(cfg.LogFilePath)
		if err != nil {
			logger.Fatal("create log file at %s failed: %v", cfg.LogFilePath, err)
		}
		logger.Info("will change logger output file to %s", cfg.LogFilePath)
		logger.SetOutput(logFile)
		defer logFile.Close()
	}

	if flag.NArg() != 1 {
		fmt.Printf("need one and only one argument to spectify action.\n\n")
		flag.Usage()
		return
	}

	switch action := flag.Arg(0); action {
	case "import":
		importFrom()
	case "run":
		startDaemon(*cfg)
	default:
		fmt.Printf("%s is not supported action argument \n\n", action)
		flag.Usage()
	}
}
