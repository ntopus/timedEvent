package main

import (
	"context"
	"devgit.kf.com.br/comercial/fleet-management-api/application"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"sort"
)

const (
	flagContract = "contract"
	flagConfig   = "config"
)

func runTokenGenerator(cli *cli.Context) error {
	appMan := application.NewTokenApplicationManager(cli.String(flagConfig))
	appMan.RunTokenGenerator(cli.String(flagContract))
	return nil
}

func runApplication(cli *cli.Context) error {
	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	appMan := application.NewApplicationManager(cli.String(flagConfig))
	err := appMan.RunApplication(ctx)
	if err != nil {
		return err
	}
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
	app.Description = ""
	app.Copyright = "nTopus© - n Possibilities..."
	app.EnableBashCompletion = true
	app.Action = runApplication
	app.Commands = []cli.Command{
		{
			Name:    "token-generator",
			Aliases: []string{"tg"},
			Action:  runTokenGenerator,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: flagContract + ", cnt",
				},
				cli.StringFlag{
					Name:  flagConfig + ", c",
					Value: "./config/config.json",
					Usage: "Path to config file",
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
