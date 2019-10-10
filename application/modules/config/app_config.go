package config

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
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
