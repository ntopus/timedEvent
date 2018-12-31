package config

import (
	"encoding/json"
	"os"
)

func configSample() *ConfigData {
	var config ConfigData

	config.DataBase.DbName = "time-queue"
	config.DataBase.ServerHost = "timeQueue.DB.ivanmeca.com.br"
	config.DataBase.ServerPort = "9003"
	config.DataBase.ServerUser = ""
	config.DataBase.ServerPassword = ""
	return &config
}

func generateConfigFile(filename string, data *ConfigData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = file.Write(jsondata)
	if err != nil {
		return err
	}
	return nil
}

func ConfigSample(sampleFile string) error {
	return generateConfigFile(sampleFile, configSample())
}
