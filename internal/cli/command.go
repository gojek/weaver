package cli

import (
	"fmt"
	"strings"
)

type cmdAction func(c *Context) error

type Command struct {
	name        string
	action      cmdAction
	subCommands Commands
}

type Commands []*Command

func (cmd *Command) CliHandler() string {
	return cmd.name
}

func (cmd *Command) Exec(c *Context) error {
	if cmd.action != nil {
		return cmd.action(c)
	} else {
		cliHandler := strings.Split(c.Command.FullName(), " ")[0]
		for _, eachCmd := range cmd.subCommands {
			if eachCmd.CliHandler() == cliHandler {
				return eachCmd.Exec(c)
			}
		}
	}
	return fmt.Errorf("No Command not registered for :%s", c.Command.FullName())
}

func (pc *Command) RegisterCommand(cmd *Command) error {
	if pc.action == nil {
		cliHandler := cmd.CliHandler()
		for _, eachCmd := range pc.subCommands {
			if eachCmd.CliHandler() == cliHandler {
				return fmt.Errorf("Another Command Regsitered for Cli Handler: %s", cliHandler)
			}
		}
		pc.subCommands = append(pc.subCommands, cmd)
		return nil
	} else {
		return fmt.Errorf("Command Does Not Allow SubCommand Registration")
	}
}

func NewDefaultCommand(name, usage, description string, action cmdAction) *Command {
	return &Command{name: name, action: action}
}

func NewParentCommand(name, usage, description string) *Command {
	return &Command{name: name, subCommands: []*Command{}}
}
