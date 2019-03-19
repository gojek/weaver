package cli

import (
	"fmt"
	"strings"
)

type ParentCommand struct {
	cmdRegistery map[string]*Command
	*Command
}

func (pc *ParentCommand) RegisterCommand(cmd *Command) error {
	cliHandler := cmd.CliHandler()
	if _, cmdFound := pc.cmdRegistery[cliHandler]; cmdFound {
		return fmt.Errorf("Another Command Regsitered for Cli Handler: %s", cliHandler)
	}

	pc.cmdRegistery[cliHandler] = cmd
	return nil
}

func (pc *ParentCommand) Exec(c *Context) error {
	fmt.Println(c.Command.FullName())
	cliHandler := strings.Split(c.Command.FullName(), " ")[0]
	if cmd, cmdFound := pc.cmdRegistery[cliHandler]; cmdFound {
		return cmd.Exec(c)
	}
	return fmt.Errorf("No Command not registered for :%s", c.Command.FullName())
}

func NewParentCommand(name, usage, description string) *ParentCommand {
	return &ParentCommand{cmdRegistery: make(map[string]*Command), Command: &Command{name: name}}
}
