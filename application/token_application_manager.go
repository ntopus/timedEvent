package application

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type TokenApplicationManager struct {
	configPath string
	auth       authenticate.IAuthenticate
}

const contractSeparator = ","

func NewTokenApplicationManager(configPath string) *TokenApplicationManager {
	return &TokenApplicationManager{
		configPath: configPath,
	}
}

func (app *TokenApplicationManager) RunTokenGenerator(contract string) error {
	fmt.Println("Initialize token generator")
	app.initializeConfig()
	err := app.initializeAuthenticate()
	if err != nil {
		return err
	}
	fmt.Println(contract)
	if contract == "" {
		return errors.New("want contract")
	}
	contracts := strings.Split(contract, contractSeparator)
	fmt.Println(contracts)
	contractsTrimmed := []string{}
	for _, c := range contracts {
		contractsTrimmed = append(contractsTrimmed, strings.Trim(c, " "))
	}
	t, err := app.auth.CreateToken(contractsTrimmed)
	if err != nil {
		return err
	}
	fmt.Println(t)
	return nil
}

func (app *TokenApplicationManager) initializeAuthenticate() error {
	tc := token.NewTokenCreator("fleet-management-api")
	sf, err := session_file.NewSessionFile(file_config.GetConfig().TokenFolder)
	if err != nil {
		return err
	}
	app.auth = authenticate.NewAuthenticate(tc, sf)
	return nil
}

func (app *TokenApplicationManager) initializeConfig() error {
	appConfig, err := file_config.LoadConfig(app.configPath)
	if err != nil {
		return err
	}
	file_config.SetConfig(appConfig)
	return config.InitConfig()
}
