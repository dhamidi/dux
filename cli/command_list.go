package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/dhamidi/dux"
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

	longestBlueprintName := ""
	for _, blueprintName := range blueprints {
		if len(longestBlueprintName) < len(blueprintName) {
			longestBlueprintName = blueprintName
		}
	}

	for _, blueprint := range blueprints {
		cmd.listBlueprint(ctx, blueprint, len(longestBlueprintName))
	}
	return cmd, nil
}

// listBlueprint displays information about a single blueprint
func (cmd *CommandList) listBlueprint(ctx *CLI, blueprintName string, labelWidth int) {
	if labelWidth < 20 {
		labelWidth = 20
	}
	entryFormat := fmt.Sprintf("%%-%ds", labelWidth)
	blueprint := new(dux.Blueprint)
	description := ""
	err := ctx.app.Store.Get(blueprintName, blueprint)
	if err == nil {
		description = blueprint.Description
	}

	fmt.Fprintf(ctx.out, entryFormat, blueprintName)
	if len(description) > 0 {
		fmt.Fprintf(ctx.out, " # %s", description)
	}
	fmt.Fprintf(ctx.out, "\n")
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
