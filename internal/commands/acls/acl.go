package acls

import (
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/etcd"
	"github.com/gojektech/weaver/internal/cli"
	"github.com/gojektech/weaver/internal/commands"
	"github.com/gojektech/weaver/pkg/logger"
	"os"
)

const (
	aclsCmdName        = "acls"
	aclsCmdUsage       = "ACLs - Perform CRUD Operations"
	aclsCmdDescription = "ACLs - Perform CRUD Operations"
)

var weaverACLSCmd = cli.NewParentCommandWithAction(aclsCmdName, aclsCmdUsage, aclsCmdDescription, setupACLs)

func setupACLs(c *cli.Context) error {
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
	weaverACLSCmd.SetFlag(commands.ETCDFLAG)
	weaverACLSCmd.SetFlag(commands.NAMESPACEFLAG)
	cli.RegisterAsBaseCommand(weaverACLSCmd)
}
