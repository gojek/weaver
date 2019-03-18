package cli

import ()

type ParentCommand struct {
	cmdRegistery map[string]*Command
	*Command
}

func (pc *ParentCommand) RegisterCommand(cmd *Command) error {
	pc.cmdRegistery[cmd.CliHandler()] = cmd
	return nil
}

func NewParentCommand(name, usage, description string) *ParentCommand {
	return &ParentCommand{cmdRegistery: make(map[string]*Command), Command: &Command{name: name}}
}
