package cli_test

import (
	"fmt"
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

func TestParentCommandShouldRegisterCommand(t *testing.T) {
	ts := setupTestParentCommand()
	cmd := cli.NewDefaultCommand("test", "usage", "description", func(c *cli.Context) error { return nil })
	assert.NoError(t, ts.cmd.RegisterCommand(cmd))
}

func TestParentCommandShouldNotAllowMoreThanOneCommandPerCliHandler(t *testing.T) {
	cliHandler := "test"
	ts := setupTestParentCommand()
	cmd := cli.NewDefaultCommand(cliHandler, "usage", "description", func(c *cli.Context) error { return nil })
	assert.NoError(t, ts.cmd.RegisterCommand(cmd))
	err := ts.cmd.RegisterCommand(cmd)
	assert.Error(t, err)
	assert.Equal(t, err, fmt.Errorf("Another Command Regsitered for Cli Handler: %s", cliHandler))
}

func setupTestParentCommand() *testParentCommandSetup {
	ts := &testParentCommandSetup{name: "parent", usage: "parent usage", description: "parent description"}
	ts.cmd = cli.NewParentCommand(ts.name, ts.usage, ts.description)
	return ts
}
