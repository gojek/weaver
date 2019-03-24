package cli

import (
	"github.com/gojektech/weaver/etcd"
	baseCli "gopkg.in/urfave/cli.v1"
)

type Context struct {
	RouteLoader *etcd.RouteLoader
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
