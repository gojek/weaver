package cli

import (
	"fmt"
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/pkg/logger"
	cliParser "gopkg.in/urfave/cli.v1"
	"os"
	"strings"
)

var commandMap = make(map[string]Command)

type command struct {
	cliName     string
	description string
	usage       string
	aliases     []string
	getCommand  func() cliParser.Command
	action      func(c *cliParser.Context) error
}

type Command interface {
	GetCommand() cliParser.Command
	Exec(c *cliParser.Context) error
}

func (cmd command) GetCommand() cliParser.Command {
	return cmd.getCommand()
}

func (cmd command) Exec(c *cliParser.Context) error {
	return cmd.action(c)
}

func registerCommand(cliName string, cmd Command) {
	if cmdToExecute := commandMap[cliName]; cmdToExecute != nil {
		panic(fmt.Sprintf("Command for cli %s already registered by: %s", cliName, cmdToExecute))
	}
	commandMap[cliName] = cmd
}

func GetCommands() []cliParser.Command {
	cmds := []cliParser.Command{}
	for _, registerCmd := range commandMap {
		cmds = append(cmds, registerCmd.GetCommand())
	}
	return cmds
}

func ExecCommand(c *cliParser.Context) error {
	cliName := getCliName(c.Command)

	if cmdToExecute := commandMap[cliName]; cmdToExecute != nil {
		setUp(c)
		return cmdToExecute.Exec(c)
	}
	return fmt.Errorf("Error executing command. Command not registered: %s", c.Command.Name)
}

func getCliName(cmd cliParser.Command) string {
	return strings.Split(cmd.FullName(), " ")[0]
}

func setUp(c *cliParser.Context) {
	os.Setenv("LOGGER_LEVEL", c.GlobalString("verbose"))
	config.Load()
	logger.SetupLogger()
}
