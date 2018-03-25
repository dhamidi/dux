package cli

import (
	"flag"
	"fmt"
	"io"
)

// CommandList is a CLI command
type CommandList struct {
	*parentCommand
}

// NewCommandList creates a new, empty instance of this command.
func NewCommandList() *CommandList {
	return &CommandList{
		parentCommand: new(parentCommand),
	}
}

// Exec implements Command
func (cmd *CommandList) Exec(ctx *CLI, args []string) (Command, error) {
	blueprints, err := ctx.app.Store.List("*")
	if err != nil {
		return cmd, err
	}

	for _, blueprint := range blueprints {
		fmt.Fprintf(ctx.out, "%s\n", blueprint)
	}
	return cmd, nil
}

// Options implements Command
func (cmd *CommandList) Options() *flag.FlagSet {
	return nil
}

// Description implements HasDescription
func (cmd *CommandList) Description() string { return `List available blueprints` }

// ShowUsage implements HasUsage
func (cmd *CommandList) ShowUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s list\n\n", cmd.CommandPath())
	fmt.Fprintf(out, "List available blueprints\n")
	fmt.Fprintf(out, "\n")
}
