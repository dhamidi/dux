package dux

import (
	"fmt"
	"io"
)

// Command describes an action that can be executed by a user of Dux
type Command interface {
	Execute(ctx *Context, args []string) error
	CommandName() string
	CommandDescription() string
}

// Application exposes a set of commands to a user.
type Application struct {
	name               string
	context            *Context
	commands           map[string]Command
	longestCommandName string
}

// NewApplication creates a new application with the given context and no registered commands.
func NewApplication(name string, ctx *Context) *Application {
	return &Application{
		name:               name,
		context:            ctx,
		commands:           map[string]Command{},
		longestCommandName: "",
	}
}

// DefineCommand registers a new command constructor for the given command identifier.
func (app *Application) DefineCommand(command Command) *Application {
	name := command.CommandName()
	app.commands[name] = command
	if len(name) > len(app.longestCommandName) {
		app.longestCommandName = name
	}
	return app
}

// RunCommand runs a new command based on the given command name and passes any arguments to the command itself.
func (app *Application) RunCommand(name string, args []string) error {
	command := app.commands[name]
	if command == nil {
		return fmt.Errorf("No such command: %q", name)
	}

	return command.Execute(app.context, args)
}

// Usage describes all available commands in the application
func (app *Application) Usage(w io.Writer) error {
	fmt.Fprintf(w, `Usage: %s COMMAND

A template-based code generator for projects in any language.

Commands:
`, app.name)
	commandLineFormat := fmt.Sprintf("  %%%ds   %%-s\n", len(app.longestCommandName))
	for commandName, command := range app.commands {
		fmt.Fprintf(w, commandLineFormat, commandName, command.CommandDescription())
	}
	fmt.Fprintf(w, "\n")
	return nil
}
