package application

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"github.com/ivanmeca/timedEvent/application/modules/queue_publisher"
	"github.com/ivanmeca/timedEvent/application/modules/scheduler"
	"github.com/ivanmeca/timedEvent/application/modules/server"
	"github.com/pkg/errors"
	"os"
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
	err := app.verifyConfig(app.configPath)
	if err != nil {
		println("Erro:", err.Error())
		os.Exit(1)
	}
	app.initializeLogger()
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

func (app *ApplicationManager) initializeLogger() {
	AppLogger := logger.GetLogger()
	LogLevel := config.GetConfig().LogLevel
	AppLogger.SetLogLevel(LogLevel)
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

func (app *ApplicationManager) verifyConfig(configFile string) error {
	if len(configFile) < 3 {
		return errors.New("Could not get config file")
	}
	config.InitConfig(configFile)
	return nil
}
