package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

var strConfig ConfigData

func InitConfigOld(filename string) {
	err := config.Load(file.NewSource(
		file.WithPath(filename),
	))
	if err != nil {
		panic(err.Error())
	}
	err = config.Scan(&strConfig)
	if err != nil {
		panic(err.Error())
	}
}

func GetConfig() *ConfigData {
	return &strConfig
}
