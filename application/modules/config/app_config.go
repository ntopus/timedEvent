package config

import (
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/file"
	"log"
)

var strConfig ConfigData

func InitConfig(filename string) {
	err := config.Load(file.NewSource(
		file.WithPath(filename),
	))
	if err != nil {
		log.Fatal("could not loud config file")
	}
	err = config.Scan(&strConfig)
	if err != nil {
		log.Fatal("could not scan config file")
	}
}

func GetConfig() *ConfigData {
	return &strConfig
}
