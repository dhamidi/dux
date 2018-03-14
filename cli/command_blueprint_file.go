package cli

import (
	"flag"
	"fmt"

	"github.com/dhamidi/dux"
)

// CommandBlueprintFile is a CLI command for rendering a blueprint.
type CommandBlueprintFile struct {
	BlueprintName string
	TemplateName  string
	FileName      string
}

// NewCommandBlueprintFile creates a new, empty instance of this command.
func NewCommandBlueprintFile() *CommandBlueprintFile {
	return &CommandBlueprintFile{}
}

// Exec implements Command
func (cmd *CommandBlueprintFile) Exec(ctx *CLI, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]
	args = args[1:]
	if len(args) == 0 {
		return fmt.Errorf("No file name provided")
	}
	cmd.FileName = args[0]

	args = args[1:]
	if len(args) == 0 {
		return fmt.Errorf("No template name provided")
	}
	cmd.TemplateName = args[0]

	return ctx.app.Execute(&dux.DefineBlueprintFile{
		BlueprintName: cmd.BlueprintName,
		TemplateName:  cmd.TemplateName,
		FileName:      cmd.FileName,
	})
}

// Options implements Command
func (cmd *CommandBlueprintFile) Options() *flag.FlagSet { return nil }
