package main

import (
	"context"
	"errors"
	"github.com/ivanmeca/timedEvent/application"
	"github.com/ivanmeca/timedEvent/config"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"sort"
)

const (
	flagConfig = "config"
)

func verifyConfig(cli *cli.Context) error {
	configFile := cli.String("config")
	if len(configFile) < 3 {
		return errors.New("Could not get config file")
	}
	config.InitConfig(configFile)
	return nil
}

func runGenerateConfigSample(cli *cli.Context) error {
	configSampleFile := cli.String(flagConfig)
	err := config.ConfigSample(configSampleFile)
	if err != nil {
		log.Fatal("Could not get config sample")
		return err
	}
	return nil
}

func runApplication(cli *cli.Context) error {
	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	err := verifyConfig(cli)
	if err != nil {
		log.Fatal("Could not loud config file")
		return err
	}
	application.RunMainApplication(ctx)
	defer cancel()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return nil
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagConfig + ", c",
			Value: "./config/config.json",
			Usage: "Path to config file",
		},
	}
	app.Version = Version + "(" + GitCommit + ")"
	app.Name = ApplicationName
	app.Usage = ""
	app.Description = "Service for support timed queue"
	app.EnableBashCompletion = true
	app.Action = runApplication
	app.Commands = []cli.Command{
		{
			Name:    "config-sample",
			Usage:   "generate config application file sample",
			Aliases: []string{"cfg-sample"},
			Action:  runGenerateConfigSample,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagConfig + ", c",
					Value: "./config/config.json",
					Usage: "Path to config sample file",
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		panic(err.Error())
	}
}
