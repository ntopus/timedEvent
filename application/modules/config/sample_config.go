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

	config.DataBase.DbName = "timed-event"
	config.DataBase.ServerHost = "127.0.0.1"
	config.DataBase.ServerPort = "9003"
	config.DataBase.ServerUser = ""
	config.DataBase.ServerPassword = ""

	var cqueueconf ConfigQueue
	cqueueconf.ServerHost = "127.0.0.1"
	cqueueconf.ServerPort = "5672"
	cqueueconf.ServerUser = "dummy_user"
	cqueueconf.ServerPassword = "dummy_pass"
	cqueueconf.QueueName = "throwAt"

	config.PublishQueue = append(config.PublishQueue, cqueueconf)
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
	return generateConfigFile("./config-sample.json", configSample())
}
