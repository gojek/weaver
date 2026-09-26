package acls

import (
	"github.com/gojektech/weaver/internal/cli"
	"github.com/gojektech/weaver/internal/views"
	"github.com/gojektech/weaver/pkg/logger"
)

const (
	listCmdName        = "list"
	listCmdUsage       = "List Weaver ACLS in ETCD Under a Namespace"
	listCmdDescription = "List Weaver ACLS in ETCD Under a Namespace"
)

var aclListCmd = cli.NewDefaultCommand(listCmdName, listCmdUsage, listCmdDescription, listACL)

func listACL(c *cli.Context) error {
	acls, err := c.RouteLoader.ListAll()
	if err != nil {
		logger.Fatalf("Error while listing acls: %s", err)
	}

	type aclInfo struct {
		ID        string `json:"ACL ID"`
		Criterion string `json:"Criterion"`
	}

	formattedAcls := []aclInfo{}

	for _, eachACL := range acls {
		formattedAcls = append(formattedAcls, aclInfo{eachACL.ID, eachACL.Criterion})
	}

	views.Render(formattedAcls)
	return nil
}

func init() {
	weaverACLSCmd.RegisterCommand(aclListCmd)
}
