package cli

import (
	"fmt"
	"github.com/gojektech/weaver/etcd"
	baseCli "gopkg.in/urfave/cli.v1"
)

var registeredCommands = Commands{}

type Context struct {
	RouteLoader *etcd.RouteLoader
	*baseCli.Context
}

func RegisterAsBaseCommand(cmd *Command) error {
	cliHandler := cmd.CliHandler()
	for _, eachCmd := range registeredCommands {
		if eachCmd.CliHandler() == cliHandler {
			return fmt.Errorf("Another Command Regsitered for Cli Handler: %s", cliHandler)
		}
	}
	registeredCommands = append(registeredCommands, cmd)
	return nil
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
