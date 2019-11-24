package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)


// Parse parse ini config and check it.
func Parse(configPath string) (*Config, error) {
	INIConfig, err := ini.Load(configPath)
	var mc Config
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("general").MapTo(&mc)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("admin").MapTo(&mc.Admin)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("moderation ").MapTo(&mc.Moderation)
	if err != nil {
		return nil, err
	}
	err = INIConfig.Section("server").MapTo(&mc.Server)
	fmt.Println(mc)
	return &mc, err
}
