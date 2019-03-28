package server

import (
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/etcd"
	"github.com/gojektech/weaver/internal/cli"
	"github.com/gojektech/weaver/internal/commands"
	"github.com/gojektech/weaver/pkg/logger"
	"os"
)

const (
	serverCmdName        = "server"
	serverCmdUsage       = "Weaver - Run Server"
	serverCmdDescription = "Weaver - Run Server"
)

var weaverServerCmd = cli.NewParentCommandWithAction(serverCmdName, serverCmdUsage, serverCmdDescription, setupServer)

func setupServer(c *cli.Context) error {
	os.Setenv("ETCD_ENDPOINTS", c.GlobalString("etcd-host"))
	os.Setenv("ETCD_KEY_PREFIX", c.GlobalString("namespace"))
	config.Load()
	rl, err := etcd.NewRouteLoader()

	if err != nil {
		logger.Fatalf("Couldn't create route loader: %s", err)
		os.Exit(1)
	}

	c.RouteLoader = rl
	return nil
}

func init() {
	weaverServerCmd.SetFlag(commands.ETCDFLAG)
	weaverServerCmd.SetFlag(commands.NAMESPACEFLAG)
	cli.RegisterAsBaseCommand(weaverServerCmd)
}
