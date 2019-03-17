package cli_test

import (
	"github.com/gojektech/weaver/internal/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSetup struct {
	name, usage, description string
	isCalled                 bool
	cmd                      *cli.Command
}

func TestDefaultCommandInitialization(t *testing.T) {
	ts := setup()
	assert.NotNil(t, ts.cmd)
}

func TestCommandShouldHaveCliHandlerName(t *testing.T) {
	ts := setup()
	assert.Equal(t, ts.cmd.CliHandler(), ts.name)
}

func TestCommandShouldExecuteSpecifiedAction(t *testing.T) {
	ts := setup()
	ts.cmd.Exec(&cli.Context{})
	assert.True(t, ts.isCalled)
}

func setup() *testSetup {
	setup := &testSetup{name: "test", usage: "usage", description: "description", isCalled: false}
	action := func(c *cli.Context) error { setup.isCalled = true; return nil }
	setup.cmd = cli.NewDefaultCommand(setup.name, setup.usage, setup.description, action)
	return setup
}
