package dux

// Command describes the arguments for an action that should be performed by dux.
type Command interface {
	CommandName() string
}

// CommandHandler is implemented by objects that can execute actions based on commands.
type CommandHandler interface {
	Execute(command Command) error
}

// CommandHandlerFunc implements CommandHandler by calling the wrapped function
type CommandHandlerFunc func(Command) error

// Execute calls the underlying function with command
func (f CommandHandlerFunc) Execute(command Command) error {
	return f(command)
}
