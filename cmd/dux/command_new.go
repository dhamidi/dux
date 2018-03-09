package main

import "github.com/dhamidi/dux"

// CommandNew is a command for creating new files from a blueprint.
type CommandNew struct{}

// CommandName implements dux.Command
func (c *CommandNew) CommandName() string { return "new" }

// CommandDescription implements dux.Command
func (c *CommandNew) CommandDescription() string { return `Generate files from a blueprint.` }

// Execute renders templates from the blueprint
func (c *CommandNew) Execute(ctx *dux.Context, args []string) error {
	blueprintName := args[0]
	blueprint, err := ctx.LoadBlueprint(blueprintName)
	if err != nil {
		return err
	}
	if err := blueprint.ParseArgs(args[1:]); err != nil {
		return err
	}

	if err := blueprint.GatherData(); err != nil {
		return err
	}

	result := blueprint.Render(ctx)
	if result.HasError() {
		return result
	}
	if err := blueprint.RenderEdits(ctx); err != nil {
		return err
	}

	return blueprint.CopyFilesToDestination(ctx, result)
}
