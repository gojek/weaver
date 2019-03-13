package cli

import (
	"github.com/gojektech/weaver/etcd"
	views "github.com/gojektech/weaver/internal/views/acls"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	aclListCliName     = "list"
	aclListDescription = "List ALL Weaver ACL in a Namespace"
	aclListUsage       = "List ALL Weaver ACL in a Namespace"
)

var (
	aclListAlises = []string{"l"}

	aclList = command{
		cliName:     aclListCliName,
		description: aclListDescription,
		usage:       aclListUsage,
		aliases:     aclListAlises,

		getCommand: func() cli.Command {
			return cli.Command{
				Name:        aclListCliName,
				Description: aclListDescription,
				Usage:       aclListUsage,
				Aliases:     aclListAlises,
				Action:      ExecCommand,
			}
		},

		action: func(c *cli.Context) error {
			rl, err := etcd.NewRouteLoader()
			if err != nil {
				return err
			}
			acls, err := rl.ListAll()
			if err != nil {
				return err
			}
			views.RenderList(acls)
			return nil

		},
	}
)

func init() {
	registerACLCommand(aclListCliName, aclList)
}
