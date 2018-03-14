package cli

import (
	"flag"
	"fmt"

	"github.com/dhamidi/dux"
)

// CommandNew is a CLI command for rendering a blueprint.
type CommandNew struct {
	BlueprintName string
	Destination   string
}

// NewCommandNew creates a new, empty instance of this command.
func NewCommandNew() *CommandNew {
	return &CommandNew{}
}

// Exec implements Command
func (cmd *CommandNew) Exec(ctx *CLI, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]

	return ctx.app.Execute(&dux.RenderBlueprint{
		Name:        cmd.BlueprintName,
		Destination: ".dux",
		Data:        map[string]interface{}{},
	})
}

// Options implements Command
func (cmd *CommandNew) Options() *flag.FlagSet {
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.StringVar(&cmd.Destination, "destination", "", "Directory into which to render the blueprint")
	return flags
}
