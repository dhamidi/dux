package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/dhamidi/dux"
)

// CommandBlueprintCreate is a CLI command
type CommandBlueprintCreate struct {
	*parentCommand
}

// NewCommandBlueprintCreate creates a new, empty instance of this command.
func NewCommandBlueprintCreate() *CommandBlueprintCreate {
	return &CommandBlueprintCreate{
		parentCommand: new(parentCommand),
	}
}

// Exec implements Command
func (cmd *CommandBlueprintCreate) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) == 0 {
		return cmd, fmt.Errorf("No blueprint name given")
	}

	err := ctx.app.Execute(&dux.CreateBlueprint{
		Name: args[0],
	})
	if err != nil {
		return cmd, err
	}

	return cmd, nil
}

// Options implements Command
func (cmd *CommandBlueprintCreate) Options() *flag.FlagSet {
	return nil
}

// Description implements HasDescription
func (cmd *CommandBlueprintCreate) Description() string { return `Create a new blueprint` }

// ShowUsage implements HasUsage
func (cmd *CommandBlueprintCreate) ShowUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s create NAME\n\n", cmd.CommandPath())
	fmt.Fprintf(out, "Create a new blueprint called NAME")
	fmt.Fprintf(out, "\n")
}
