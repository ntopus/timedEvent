package application

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/queue_publisher"
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
	cancelServer := app.initializeServer()
	cancelScheduler := app.initializeScheduler()
	queue_publisher.QueuePublisher()
	go func() {
		<-ctx.Done()
		cancelServer()
		cancelScheduler()
	}()
	return nil
}

func (app *ApplicationManager) initializeServer() context.CancelFunc {
	s := server.NewHttpServer(strconv.Itoa(config.GetConfig().ServerPort), false)
	ctxServer := context.Background()
	ctxServer, cancelServer := context.WithCancel(ctxServer)
	s.RunServer(ctxServer)
	return cancelServer
}

func (app *ApplicationManager) initializeScheduler() context.CancelFunc {
	s := scheduler.NewScheduler(config.GetConfig().PoolTime, config.GetConfig().ControlTime, config.GetConfig().ExpirationTime)
	ctxServer := context.Background()
	ctxServer, cancelServer := context.WithCancel(ctxServer)
	s.Run(ctxServer)
	return cancelServer
}
