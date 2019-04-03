package main

import (
	"context"
	"errors"
	"github.com/ivanmeca/timedEvent/application"
	"github.com/ivanmeca/timedEvent/config"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"sort"
)

const flagConfig = "config"
const configFilePathDefault = "./config/config.json"

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
	case <-appMan.Done():
		return nil
	}
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
	gin.SetMode(gin.ReleaseMode)
	if VersionPrerelease != "" {
		app.Version += " - " + VersionPrerelease
		gin.SetMode(gin.DebugMode)
	}
	app.Name = ApplicationName
	app.Usage = ""
	app.Description = ""
	app.Copyright = "nTopus© - n Possibilities..."
	app.EnableBashCompletion = true
	app.Action = runApplication
	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
