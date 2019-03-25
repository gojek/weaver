package cli

import (
	"fmt"
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/etcd"
	"github.com/gojektech/weaver/pkg/logger"
	baseCli "gopkg.in/urfave/cli.v1"
	"os"
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

func GetBaseCommands() []baseCli.Command {
	baseCliCommands := []baseCli.Command{}
	for _, eachCmd := range registeredCommands {
		baseCmd := baseCli.Command{
			Name:        eachCmd.name,
			Usage:       eachCmd.usage,
			Description: eachCmd.description,
			Flags:       eachCmd.flags,
			Action: func(ctx *baseCli.Context) error {
				c := &Context{Context: ctx}
				setup(c)
				return eachCmd.Exec(c)
			},
		}
		baseCliCommands = append(baseCliCommands, baseCmd)
	}
	return baseCliCommands
}

func setup(c *Context) {
	os.Setenv("LOGGER_LEVEL", c.GlobalString("verbose"))
	config.Load()
	logger.SetupLogger()
}

func NewApp() *baseCli.App {
	app := baseCli.NewApp()
	app.Flags = []baseCli.Flag{
		NewStringFlag("verbose", "Error", "Verbosity of log level, ex: debug, info, warn, error, fatal, panic", "LOGGER_LEVEL"),
	}
	return app
}
