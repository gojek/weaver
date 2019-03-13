package acls

import (
	"github.com/gojektech/weaver"
	"github.com/jedib0t/go-pretty/table"
	"os"
)

func RenderList(acls []*weaver.ACL) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Criterion"})
	for id, acl := range acls {
		t.AppendRow([]interface{}{id, acl.ID, acl.Criterion})
	}
	t.Render()
}
