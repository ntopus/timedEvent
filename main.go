package main

import (
	"context"
	"errors"
	"github.com/ivanmeca/timedEvent/application"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"sort"
)

const flagConfig = "config"
const configFilePathDefault = "./config/config.json"

func runApplication(cli *cli.Context) error {
	err := verifyConfig(cli)
	if err != nil {
		println("Erro:", err.Error())
		os.Exit(1)
	}
	AppLogger := logger.GetLogger()
	LogLevel := config.GetConfig().LogLevel
	AppLogger.SetLogLevel(LogLevel)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	appMan := application.NewApplicationManager(cli.String(flagConfig))
	err = appMan.RunApplication(ctx)
	if err != nil {
		return err
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-quit:
		cancel()
		return nil
	}
}

func verifyConfig(cli *cli.Context) error {
	configFile := cli.String("config")
	if len(configFile) < 3 {
		return errors.New("Could not get config file")
	}
	config.InitConfig(configFile)
	return nil
}

func generateConfig(cli *cli.Context) error {
	err := config.ConfigSample()
	if err != nil {
		log.Fatal("Could not get config sample")
		return err
	}
	return nil
	//
	//qUser, ok := cfg.GetConfigParam(config.QueueServerUser, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get queue user")
	//	os.Exit(1)
	//}
	//qPass, ok := cfg.GetConfigParam(config.QueueServerPassword, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get queue password")
	//	os.Exit(1)
	//}
	//qHost, ok := cfg.GetConfigParam(config.QueueServerHost, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get queue host")
	//	os.Exit(1)
	//}
	//qPort, ok := cfg.GetConfigParam(config.QueueServerPort, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get queue port")
	//	os.Exit(1)
	//}
	//qName, ok := cfg.GetConfigParam(config.QueueName, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get queue name")
	//	os.Exit(1)
	//}
	//qErrorName, ok := cfg.GetConfigParam(config.QueueErrorName, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get Error queue name")
	//	os.Exit(1)
	//}
	//qNotifyName, ok := cfg.GetConfigParam(config.QueueNotifyName, "").(string)
	//if !ok {
	//	AppLogger.CriticalPrintln("Could not get Notify queue name")
	//	os.Exit(1)
	//}
	//qPortInt, err := strconv.Atoi(qPort)
	//if err != nil {
	//	AppLogger.CriticalPrintln("Could not convert queue port value")
	//	os.Exit(1)
	//}
	//
	//qpParams := queue_repository.NewQueueRepositoryParams(qUser, qPass, qHost, qPortInt)
	//qRepo, err := queue_repository.NewQueueRepository(qpParams)
	//if err != nil {
	//	AppLogger.CriticalPrintln("Error on queue start:", err)
	//	os.Exit(1)
	//}
	//qp, err := qRepo.QueueDeclare(queue.NewQueueParams(qName), false)
	//if err != nil {
	//	AppLogger.CriticalPrintln("Error on Error queue start:", err)
	//	os.Exit(1)
	//}

}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagConfig + ", c",
			Value: configFilePathDefault,
			Usage: "Path to config file",
		},
	}
	app.Version = Version + "(" + GitCommit + ")"
	app.Name = ApplicationName
	app.Usage = ""
	app.Description = ""
	app.Copyright = "nTopus© - n Possibilities..."
	app.EnableBashCompletion = true
	app.Action = runApplication
	app.Commands = []cli.Command{
		{
			Name:    "config-sample",
			Aliases: []string{"cs"},
			Action:  generateConfig,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
