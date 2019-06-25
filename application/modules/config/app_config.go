package config

import (
	"encoding/json"
	"io/ioutil"
)

type AppConfig struct {
	configFile string
	configData map[string]interface{}
}

func NewServerConfigByFile(filename string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		appConfig := configSample(filename)
		err := generateConfigFile(appConfig)
		if err != nil {
			return nil, err
		}
		return appConfig, nil
	}

	appConfig := AppConfig{configFile: filename, configData: make(map[string]interface{})}

	err = json.Unmarshal([]byte(data), &appConfig.configData)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}

func (sc *AppConfig) GetConfigParam(attr string, defaultValue interface{}) interface{} {
	if v, i := sc.configData[attr]; i {
		return v
	}
	return defaultValue
}

func (sc *AppConfig) SetConfigParam(attr string, Value interface{}) {
	sc.configData[attr] = Value
}

func (sc *AppConfig) SaveConfig() error {

	data, err := json.Marshal(sc.configData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(sc.configFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
