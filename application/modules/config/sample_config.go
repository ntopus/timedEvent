package config

import (
	"encoding/json"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"os"
)

func configSample() *ConfigData {
	var config ConfigData

	config.LogLevel = logger.LogDebug
	config.ControlTime = 100
	config.PoolTime = 10000
	config.ExpirationTime = 2592000000
	config.ServerPort = 9010

	config.DataBase.DbName = "testDb"
	config.DataBase.ServerHost = "http://localhost"
	config.DataBase.ServerPort = "8529"
	config.DataBase.ServerUser = "root"
	config.DataBase.ServerPassword = "rootpass"

	var pqueueconf ConfigQueue
	pqueueconf.ServerHost = "localhost"
	pqueueconf.ServerVHost = "/timed"
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

func ConfigSample(fileName string) error {
	return generateConfigFile(fileName, configSample())
}
