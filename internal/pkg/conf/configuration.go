package conf

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

// Configure save all config for this project
type Configure struct {
	*ini.File
}

// Load config for isso
func Load(confPath string) (Configure, error) {
	log.Infof("Begin to load config from %s", confPath)
	conf, err := ini.Load(confPath)
	if err != nil {
		log.Fatalf("Fail to read default config file: %v", err)
	}
	err = verifyConf(conf)
	if err != nil {
		log.Fatalf("Invalid conf:%v", err)
	}
	return Configure{conf}, nil
}

//TODO: verification configuration.
var (
	sections []string
)

func verifyConf(conf *ini.File) error {
	return nil
}
