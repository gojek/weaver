package cli_test

import (
	"github.com/gojektech/weaver/internal/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testParentCommandSetup struct {
	name, usage, description string
	cmd                      *cli.ParentCommand
}

func TestParentCommandInitialization(t *testing.T) {
	ts := setupTestParentCommand()
	assert.NotNil(t, ts.cmd)
}

func setupTestParentCommand() *testParentCommandSetup {
	ts := &testParentCommandSetup{name: "parent", usage: "parent usage", description: "parent description"}
	ts.cmd = cli.NewParentCommand(ts.name, ts.usage, ts.description)
	return ts
}
