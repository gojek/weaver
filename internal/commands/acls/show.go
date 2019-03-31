package acls

import (
	"github.com/gojektech/weaver"
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
	acls, err := c.RouteLoader.ListAll()
	if err != nil {
		logger.Fatalf("Error while showing acls: %s", err)
	}

	aclToShow := []weaver.ACL{}

	for _, eachACL := range acls {
		if eachACL.ID == aclID {
			aclToShow = append(aclToShow, *eachACL)
		}
	}

	views.Render(aclToShow)
	return nil
}

func init() {
	aclShowCmd.SetFlag(cli.NewStringFlag("id", "", "ID Of the ACL", ""))
	weaverACLSCmd.RegisterCommand(aclShowCmd)
}
