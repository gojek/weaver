package cli

type ParentCommand struct {
	cmdRegistery map[string]*Command
	*Command
}

func NewParentCommand(name, usage, description string) *ParentCommand {
	return &ParentCommand{Command: &Command{name: name}}
}
