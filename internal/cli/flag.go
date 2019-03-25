package cli

import (
	baseCli "gopkg.in/urfave/cli.v1"
)

func NewStringFlag(name, value, usage, env string) baseCli.Flag {
	return baseCli.StringFlag{Name: name, Value: value, Usage: usage, EnvVar: env}
}
