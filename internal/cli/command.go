package cli

type cmdAction func(c *Context) error

type Command struct {
	name   string
	action cmdAction
}

func (cmd *Command) CliHandler() string {
	return cmd.name
}

func (cmd *Command) Exec(c *Context) error {
	return cmd.action(c)
}

func NewDefaultCommand(name, usage, description string, action cmdAction) *Command {
	return &Command{name, action}
}
