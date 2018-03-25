package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/dhamidi/dux"
)

// CommandBlueprintDescribe is a CLI command
type CommandBlueprintDescribe struct {
	*parentCommand
}

// NewCommandBlueprintDescribe creates a new, empty instance of this command.
func NewCommandBlueprintDescribe() *CommandBlueprintDescribe {
	return &CommandBlueprintDescribe{
		parentCommand: new(parentCommand),
	}
}

// Exec implements Command
func (cmd *CommandBlueprintDescribe) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) < 2 {
		return cmd, fmt.Errorf("missing arguments")
	}

	blueprintName := args[0]
	description := args[1]

	err := ctx.app.Execute(&dux.DescribeBlueprint{
		BlueprintName: blueprintName,
		Description:   description,
	})

	return cmd, err
}

// Options implements Command
func (cmd *CommandBlueprintDescribe) Options() *flag.FlagSet {
	return nil
}

// Description implements HasDescription
func (cmd *CommandBlueprintDescribe) Description() string {
	return `Set the description for a blueprint`
}

// ShowUsage implements HasUsage
func (cmd *CommandBlueprintDescribe) ShowUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s describe BLUEPRINT DESCRIPTION\n\n", cmd.CommandPath())
	fmt.Fprintf(out, "Set the description for BLUEPRINT to DESCRIPTION\n")
	fmt.Fprintf(out, "\n")
}
