package acls

import (
	baseCli "gopkg.in/urfave/cli.v1"

	"github.com/gojektech/weaver/internal/cli"
	"github.com/gojektech/weaver/internal/views"
	"github.com/gojektech/weaver/pkg/logger"
)

const (
	showCmdName        = "show"
	showCmdUsage       = "Show Weaver ACLS Given ACL ID"
	showCmdDescription = "Show Weaver ACLS Given ACL ID"
)

var aclShowCmd = cli.NewDefaultCommand(showCmdName, showCmdUsage, showCmdDescription, showACL)

func showACL(c *cli.Context) error {
	aclID := c.String("id")
	if aclID == "" {
		baseCli.ShowSubcommandHelp(c.Context)
		return nil
	}
	acls, err := c.RouteLoader.ListAll()
	if err != nil {
		logger.Fatalf("Error while showing acls: %s", err)
	}

	for _, eachACL := range acls {
		if eachACL.ID == aclID {
			views.Render(eachACL)
			return nil
		}

	}
	logger.Fatalf("ACL with id: %s not found", aclID)
	return nil
}

func init() {
	aclShowCmd.SetFlag(cli.NewStringFlag("id", "", "ID Of the ACL", ""))
	weaverACLSCmd.RegisterCommand(aclShowCmd)
}
