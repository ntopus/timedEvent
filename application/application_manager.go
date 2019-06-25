package application

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/scheduler"
	"github.com/ivanmeca/timedEvent/application/modules/server"
	"strconv"
)

type ApplicationManager struct {
	configPath string
}

func NewApplicationManager(configPath string) *ApplicationManager {
	return &ApplicationManager{
		configPath: configPath,
	}
}
func (app *ApplicationManager) RunApplication(ctx context.Context) error {
	fmt.Println("Initialize API")
	initializeDB()
	cancelServer := app.initializeServer()
	cancelScheduler := app.initializeScheduler()
	go func() {
		<-ctx.Done()
		cancelServer()
		cancelScheduler()
	}()
	return nil
}

func initializeDB() {
	//err := fleetDB.EnsureIndex()
	//if err != nil {
	//	panic(err.Error())
	//}
}

func (app *ApplicationManager) initializeServer() context.CancelFunc {
	//s := server.NewHttpServer(strconv.Itoa(file_config.GetConfig().Port))
	s := server.NewHttpServer(strconv.Itoa(9010))
	ctxServer := context.Background()
	ctxServer, cancelServer := context.WithCancel(ctxServer)
	s.RunServer(ctxServer)
	return cancelServer
}

func (app *ApplicationManager) initializeScheduler() context.CancelFunc {
	s := scheduler.NewScheduler(config.GetConfig().PoolTime, config.GetConfig().PoolTime)
	ctxServer := context.Background()
	ctxServer, cancelServer := context.WithCancel(ctxServer)
	s.Run(ctxServer)
	return cancelServer
}
