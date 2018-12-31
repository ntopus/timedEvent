package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

var strConfig ConfigData

func InitConfig(filename string) {
	config.Load(file.NewSource(
		file.WithPath(filename),
	))
	config.Scan(&strConfig)
}

func GetConfig() *ConfigData {
	return &strConfig
}
