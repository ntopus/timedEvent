package application

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/authenticate"
	"github.com/ivanmeca/timedEvent/application/modules/config/file_config"
	"github.com/ivanmeca/timedEvent/application/modules/server"
	"strconv"
)

type ApplicationManager struct {
	configPath string
	auth       authenticate.IAuthenticate
}

func NewApplicationManager(configPath string) *ApplicationManager {
	return &ApplicationManager{
		configPath: configPath,
	}
}
func (app *ApplicationManager) RunApplication(ctx context.Context) error {
	fmt.Println("Initialize API")
	app.initializeConfig()
	initializeDB()
	cancelServer := app.initializeServer()
	go func() {
		<-ctx.Done()
		cancelServer()
	}()
	return nil
}

func (app *ApplicationManager) initializeConfig() error {
	appConfig, err := file_config.LoadConfig(app.configPath)
	if err != nil {
		return err
	}
	file_config.SetConfig(appConfig)
	return config.InitConfig()
}

func initializeDB() {
	err := fleetDB.EnsureIndex()
	if err != nil {
		panic(err.Error())
	}
}

func (app *ApplicationManager) initializeServer() context.CancelFunc {
	s := server.NewHttpServer(strconv.Itoa(file_config.GetConfig().Port), app.auth)
	ctxServer := context.Background()
	ctxServer, cancelServer := context.WithCancel(ctxServer)
	s.RunServer(ctxServer)
	return cancelServer
}
