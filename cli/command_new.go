package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/dhamidi/dux"
)

// CommandNew is a CLI command for rendering a blueprint.
type CommandNew struct {
	BlueprintName string
	Destination   string
	DryRun        bool
}

// NewCommandNew creates a new, empty instance of this command.
func NewCommandNew() *CommandNew {
	return &CommandNew{}
}

// Exec implements Command
func (cmd *CommandNew) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) == 0 {
		return cmd, fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]

	sources := []string{}
	destinations := []string{}
	done := ctx.app.EventStore.Subscribe(cmd.collectRenderedFiles(&sources, &destinations))

	if err := ctx.app.Execute(&dux.RenderBlueprint{
		Name:        cmd.BlueprintName,
		Destination: ".dux",
		Data:        map[string]interface{}{},
	}); err != nil {
		return cmd, err
	}
	done()

	if cmd.DryRun {
		return cmd, nil
	}
	return cmd, ctx.app.Execute(&dux.Install{
		Sources:      sources,
		Destinations: destinations,
	})
}

// collectRenderedFiles listens to events emitted by RenderBlueprint to build a list of files to install.
func (cmd *CommandNew) collectRenderedFiles(sources *[]string, destinations *[]string) func(*dux.Event) {
	return func(e *dux.Event) {
		if e.Name != "template-rendered" {
			return
		}

		source := e.Payload["filename"].(string)
		*sources = append(*sources, source)
		destination := strings.Replace(source, ".dux", ".", 1)
		*destinations = append(*destinations, destination)
	}
}

// Options implements Command
func (cmd *CommandNew) Options() *flag.FlagSet {
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.StringVar(&cmd.Destination, "destination", "", "Directory into which to render the blueprint")
	flags.BoolVar(&cmd.DryRun, "dry-run", false, "Do not move generated files into current directory")
	return flags
}

// Description implements HasDescription
func (cmd *CommandNew) Description() string { return `Create new files from blueprint` }

// ShowUsage implements HasUsage
func (cmd *CommandNew) ShowUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: new [--dry-run] BLUEPRINT\n\n")
	fmt.Fprintf(out, "Options:\n")
	fmt.Fprintf(out, " --dry-run=false   Do not move generated files into current directory\n")
	fmt.Fprintf(out, "\n")
}
