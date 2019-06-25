package config

import (
	"encoding/json"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"os"
	"strconv"
)

func generateConfigFile(appConfig *AppConfig) error {

	file, err := os.Create(appConfig.configFile)
	if err != nil {
		return err
	}

	data, err := json.Marshal(appConfig.configData)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

const (
	DatabaseServerHost     = "databaseServerHost"
	DatabaseServerPort     = "databaseServerPort"
	DatabaseServerUser     = "databaseServerUser"
	DatabaseServerPassword = "databaseServerPassword"
	QueueServerHost        = "queueServerHost"
	QueueServerPort        = "queueServerPort"
	QueueServerUser        = "queueServerUser"
	QueueServerPassword    = "queueServerPassword"
	QueueName              = "queueName"
	AppToken               = "AppToken"
	LogLevel               = "LogLevel"
	LogFile                = "LogFile"
)

func configSample(filename string) *AppConfig {

	appConfig := AppConfig{configFile: filename, configData: make(map[string]interface{})}

	appConfig.configData[QueueServerHost] = "srvqueue.vehicular.module.ntopus.com.br"
	appConfig.configData[AppToken] = "123456"
	appConfig.configData[QueueServerPort] = "5672"
	appConfig.configData[QueueServerUser] = "miseravi"
	appConfig.configData[QueueServerPassword] = "trAfr@guR36a"
	appConfig.configData[QueueName] = "gsmSyncQueue"
	appConfig.configData[LogLevel] = strconv.Itoa(logger.LogNotice)
	appConfig.configData[LogFile] = ""
	return &appConfig
}
