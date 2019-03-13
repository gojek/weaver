package acls

import (
	"encoding/json"
	"github.com/gojektech/weaver"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/gojektech/weaver/pkg/shard"
	"github.com/jedib0t/go-pretty/table"
	"os"
)

func RenderShow(acl *weaver.ACL) {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendRow([]interface{}{"ID", acl.ID})
	t.AppendRow([]interface{}{"Criterion", acl.Criterion})
	t.AppendRow([]interface{}{"-"})
	t.AppendRow([]interface{}{"EndPoint Config"})
	t.AppendRow([]interface{}{"", "Matcher", acl.EndpointConfig.Matcher})
	t.AppendRow([]interface{}{"", "Shard Expression", acl.EndpointConfig.ShardExpr})
	t.AppendRow([]interface{}{"", "Shard Function", acl.EndpointConfig.ShardFunc})

	shardConfig := map[string]shard.BackendDefinition{}
	if err := json.Unmarshal(acl.EndpointConfig.ShardConfig, &shardConfig); err == nil {
		t.AppendRow([]interface{}{"-"})
		t.AppendRow([]interface{}{"Backends"})
		for key, backend := range shardConfig {
			t.AppendRow([]interface{}{"", key})
			t.AppendRow([]interface{}{"", "", backend.BackendName})
			t.AppendRow([]interface{}{"", "", backend.BackendURL})
			if backend.Timeout != nil {
				t.AppendRow([]interface{}{"", "", backend.Timeout})
			}
		}
	} else {
		logger.Debugf("Error %s", err)
	}
	t.Render()
}
