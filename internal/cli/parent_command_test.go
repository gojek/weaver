package cli_test

import (
	"fmt"
	"github.com/gojektech/weaver/internal/cli"
	"github.com/stretchr/testify/assert"
	baseCli "gopkg.in/urfave/cli.v1"
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

func TestParentCommandShouldExecuteSubCommand(t *testing.T) {
	ts := setupTestParentCommand()
	isCmdOneCalled := false
	isCmdTwoCalled := false
	cmdOne := cli.NewDefaultCommand("test-one", "usage", "description", func(c *cli.Context) error { isCmdOneCalled = true; return nil })
	cmdTwo := cli.NewDefaultCommand("test-two", "usage", "description", func(c *cli.Context) error { isCmdTwoCalled = true; return nil })
	errFromCmdOne := ts.cmd.RegisterCommand(cmdOne)
	errFromCmdTwo := ts.cmd.RegisterCommand(cmdTwo)
	ctx := &cli.Context{Context: &baseCli.Context{Command: baseCli.Command{Name: "test-one"}}}
	ts.cmd.Exec(ctx)

	assert.NoError(t, errFromCmdOne)
	assert.NoError(t, errFromCmdTwo)
	assert.True(t, isCmdOneCalled)
	assert.False(t, isCmdTwoCalled)
}

func setupTestParentCommand() *testParentCommandSetup {
	ts := &testParentCommandSetup{name: "parent", usage: "parent usage", description: "parent description"}
	ts.cmd = cli.NewParentCommand(ts.name, ts.usage, ts.description)
	return ts
}
