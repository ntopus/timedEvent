package config

import (
	"encoding/json"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"os"
)

func configSample() *ConfigData {
	var config ConfigData

	config.LogLevel = logger.LogNotice
	config.ControlTime = 1
	config.PoolTime = 5
	config.ExpirationTime = 1800
	config.ServerPort = 9010

	config.DataBase.DbName = "testDb"
	config.DataBase.ServerHost = "http://localhost"
	config.DataBase.ServerPort = "8529"
	config.DataBase.ServerUser = "testUser"
	config.DataBase.ServerPassword = "123456"

	var pqueueconf ConfigQueue
	pqueueconf.ServerHost = "127.0.0.1"
	pqueueconf.ServerVHost = "/"
	pqueueconf.ServerPort = "5672"
	pqueueconf.ServerUser = "randomUser"
	pqueueconf.ServerPassword = "randomPass"
	pqueueconf.QueueName = "throwAt"

	config.PublishQueue = append(config.PublishQueue, pqueueconf)
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

func ConfigSample() error {
	return generateConfigFile("./config.json", configSample())
}
