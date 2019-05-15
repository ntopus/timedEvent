package config

import (
	"encoding/json"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
)

var actual map[string]interface{}

func InitConfig() error {
	conf := config.NewConfig()
	src := env.NewSource()
	conf.Load(src)

	source := env.NewSource()
	c, err := source.Read()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(c.Data, &actual); err != nil {
		return err
	}
	return nil
}

func GetDatabaseHost() string {
	actualDB, ok := actual["mongodb"].(map[string]interface{})
	if !ok {
		return ""
	}
	server, ok := actualDB["server"].(string)
	if !ok {
		return ""
	}
	return server
}
