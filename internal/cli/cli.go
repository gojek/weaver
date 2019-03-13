package cli

import (
	cliParser "gopkg.in/urfave/cli.v1"
)

func NewApp() *cliParser.App {
	app := cliParser.NewApp()
	app.Flags = []cliParser.Flag{
		cliParser.StringFlag{
			Name:   "verbose",
			Value:  "Error",
			Usage:  "Verbosity of log level, ex: debug, info, warn, error, fatal, panic",
			EnvVar: "LOGGER_LEVEL",
		},
	}

	return app
}
