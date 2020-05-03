package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/kr/pretty"
	"gopkg.in/ini.v1"
	"wrong.wang/x/go-isso/logger"
)

// Parse parse ini config and check it.
func Parse(configPath string) (*Config, error) {
	INIConfig, err := ini.LoadSources(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, configPath)
	var mc Config
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("general").MapTo(&mc)
	if err != nil {
		return nil, err
	}
	DurMaxAge, err := time.ParseDuration(INIConfig.Section("general").Key("max-age").MustString("1m"))
	mc.MaxAge = int(DurMaxAge.Seconds())

	splitStringtoStrings(&mc.Host, "\n")
	splitStringtoStrings(&mc.Notify, ",")
	err = INIConfig.Section("admin").MapTo(&mc.Admin)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("moderation").MapTo(&mc.Moderation)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("server").MapTo(&mc.Server)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("guard").MapTo(&mc.Server.Guard)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("markup").MapTo(&mc.Server.Guard.Markup)
	if err != nil {
		return nil, err
	}
	splitStringtoStrings(&mc.Server.Guard.Markup.AllowedElements, ",")
	splitStringtoStrings(&mc.Server.Guard.Markup.AllowedAttributes, ",")
	err = INIConfig.Section("smtp").MapTo(&mc.SMTP)
	if err != nil {
		return nil, err
	}
	logger.Debug(fmt.Sprintf("%# v", pretty.Formatter(mc)))
	return &mc, err
}

func splitStringtoStrings(s *[]string, sep string) {
	if len(*s) == 0 {
		return
	}
	*s = strings.Split((*s)[0], sep)
}
