package file_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func LoadConfig(fileConfigPath string) (*AppConfig, error) {
	raw, err := ioutil.ReadFile(fileConfigPath)
	if err != nil {
		return nil, err
	}
	var c AppConfig
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to decode into struct, %v", err))
	}
	SetConfig(&c)
	return &c, nil
}
