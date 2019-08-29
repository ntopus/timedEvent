package main

import (
	"context"
	"github.com/ivanmeca/timedEvent/application"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"sort"
)

const flagConfig = "config"
const configFilePathDefault = "./config.json"

func runApplication(cli *cli.Context) error {
	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	appMan := application.NewApplicationManager(cli.String(flagConfig))
	err := appMan.RunApplication(ctx)
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

func generateConfig(cli *cli.Context) error {
	err := config.ConfigSample("./config.json")
	if err != nil {
		log.Fatal("Could not get config sample")
		return err
	}
	return nil
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
