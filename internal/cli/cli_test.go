package cli_test

import (
	"fmt"
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/internal/cli"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAppSetup struct {
	name, usage, description string
	cmd                      *cli.Command
}

func TestAppShouldRegisterACommand(t *testing.T) {
	ts := setupTestApp()
	err := cli.RegisterAsBaseCommand(ts.cmd)
	assert.NoError(t, err)
}

func TestAppShouldReturnErrorOnDuplicateRegistration(t *testing.T) {
	// This will throw error as previous command is also registered with same cli handler
	ts := setupTestApp()
	err := cli.RegisterAsBaseCommand(ts.cmd)

	assert.Error(t, err)
	assert.Equal(t, err, fmt.Errorf("Another Command Regsitered for Cli Handler: %s", ts.name))
}

func TestCliGetCommandsShouldGiveBaseCommands(t *testing.T) {
	ts := setupTestApp()
	baseCliCommands := cli.GetBaseCommands()

	assert.Equal(t, len(baseCliCommands), 1)
	assert.Equal(t, baseCliCommands[0].Name, ts.name)
	assert.Equal(t, baseCliCommands[0].Usage, ts.usage)
	assert.Equal(t, baseCliCommands[0].Description, ts.description)
}

func TestCliGetCommandsExecutionSHouldSetupConfigAndLogger(t *testing.T) {
	ts := setupTestApp()

	baseCliCommands := cli.GetBaseCommands()
	app := cli.NewApp()
	app.Commands = baseCliCommands
	app.Run([]string{"binary", "--verbose", "debug", ts.name})

	// config is supposed to have logger level set
	// If it setup logger, logging shouldn't panic
	assert.NotPanics(t, func() { logger.Info("Should not panic if logger is setup") })
	assert.Equal(t, config.LogLevel(), "debug")
}

func TestAppRunWithOSArgsShouldExecuteBaseCommandAction(t *testing.T) {
	isCmdActionExecuted := false
	ts := setupTestApp()
	ts.name = "exec-test"
	ts.cmd = cli.NewDefaultCommand(ts.name, ts.usage, ts.description, func(c *cli.Context) error { isCmdActionExecuted = true; return nil })
	err := cli.RegisterAsBaseCommand(ts.cmd)
	assert.NoError(t, err)

	baseCliCommands := cli.GetBaseCommands()
	app := cli.NewApp()
	app.Commands = baseCliCommands
	app.Run([]string{"binary", "--verbose", "debug", ts.name})

	assert.True(t, isCmdActionExecuted)
}

func setupTestApp() *testAppSetup {
	setup := &testAppSetup{name: "test", usage: "usage", description: "description"}
	action := func(c *cli.Context) error { return nil }
	setup.cmd = cli.NewDefaultCommand(setup.name, setup.usage, setup.description, action)
	return setup
}
