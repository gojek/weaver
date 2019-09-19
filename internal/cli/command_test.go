package cli_test

import (
	"github.com/gojektech/weaver/internal/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCommandSetup struct {
	name, usage, description string
	isCalled                 bool
	cmd                      *cli.Command
}

func TestDefaultCommandInitialization(t *testing.T) {
	ts := setupTestCommand()
	assert.NotNil(t, ts.cmd)
}

func TestCommandShouldHaveCliHandlerName(t *testing.T) {
	ts := setupTestCommand()
	assert.Equal(t, ts.cmd.CliHandler(), ts.name)
}

func TestCommandShouldExecuteSpecifiedAction(t *testing.T) {
	ts := setupTestCommand()
	ts.cmd.Exec(&cli.Context{})
	assert.True(t, ts.isCalled)
}

func TestShouldBeAbleToSetFlags(t *testing.T) {
	ts := setupTestCommand()
	flag := cli.NewStringFlag("test", "value", "usage", "env")
	assert.NotPanics(t, func() { ts.cmd.SetFlag(flag) }, "Setting a flag panics")
}

func setupTestCommand() *testCommandSetup {
	setup := &testCommandSetup{name: "test", usage: "usage", description: "description", isCalled: false}
	action := func(c *cli.Context) error { setup.isCalled = true; return nil }
	setup.cmd = cli.NewDefaultCommand(setup.name, setup.usage, setup.description, action)
	return setup
}
