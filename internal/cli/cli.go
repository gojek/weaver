package cli

import (
	baseCli "github.com/urfave/cli"
)

type Context struct {
	*baseCli.Context
}

func NewApp() *baseCli.App {
	app := baseCli.NewApp()
	app.Flags = []baseCli.Flag{
		baseCli.StringFlag{
			Name:   "verbose",
			Value:  "Error",
			Usage:  "Verbosity of log level, ex: debug, info, warn, error, fatal, panic",
			EnvVar: "LOGGER_LEVEL",
		},
	}
	return app
}
