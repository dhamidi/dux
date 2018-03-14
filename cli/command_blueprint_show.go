package cli

import (
	"flag"
	"fmt"

	"github.com/dhamidi/dux"
)

// CommandBlueprintShow is a CLI command for rendering a blueprint.
type CommandBlueprintShow struct {
	BlueprintName string
}

// NewCommandBlueprintShow creates a new, empty instance of this command.
func NewCommandBlueprintShow() *CommandBlueprintShow {
	return &CommandBlueprintShow{}
}

// Exec implements Command
func (cmd *CommandBlueprintShow) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) == 0 {
		return cmd, fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]
	blueprint := new(dux.Blueprint)
	if err := ctx.app.Store.Get(cmd.BlueprintName, blueprint); err != nil {
		return cmd, fmt.Errorf("Failed to load blueprint %q", cmd.BlueprintName)
	}

	cmd.Show(ctx, blueprint)
	return cmd, nil
}

// Options implements Command
func (cmd *CommandBlueprintShow) Options() *flag.FlagSet { return nil }

// Show displays a blueprint in the given context
func (cmd *CommandBlueprintShow) Show(ctx *CLI, blueprint *dux.Blueprint) {
	templateNames := []string{}
	done := ctx.app.EventStore.Subscribe(func(e *dux.Event) {
		if e.Name != "blueprint-template-found" {
			return
		}
		templateNames = append(templateNames, e.Payload["name"].(string))
	})
	ctx.app.Execute(&dux.ListTemplates{BlueprintName: blueprint.Name})
	done()

	fmt.Fprintf(ctx.out, "Name: %s\n", blueprint.Name)
	if len(blueprint.Files) > 0 {
		fmt.Fprintf(ctx.out, "Files:\n")
		for destination, template := range blueprint.Files {
			fmt.Fprintf(ctx.out, "  - name: %s\n", destination)
			fmt.Fprintf(ctx.out, "    template: %s\n", template)
		}
		fmt.Fprintf(ctx.out, "\n")
	}

	if len(templateNames) > 0 {
		fmt.Fprintf(ctx.out, "Templates:\n")
		for _, name := range templateNames {
			fmt.Fprintf(ctx.out, "  - %s\n", name)
		}
	}
}
