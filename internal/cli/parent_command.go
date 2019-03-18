package cli

import (
	"fmt"
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

func NewParentCommand(name, usage, description string) *ParentCommand {
	return &ParentCommand{cmdRegistery: make(map[string]*Command), Command: &Command{name: name}}
}
