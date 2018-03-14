package cli

import (
	"flag"
	"fmt"
	"io"

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

// Description implements HasDescription
func (cmd *CommandBlueprintFile) Description() string {
	return `Associate file with template in blueprint`
}

// ShowUsage implements HasUsage
func (cmd *CommandBlueprintFile) ShowUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: file BLUEPRINT FILENAME TEMPLATE\n\n")
	fmt.Fprintf(out, "Define FILENAME to be generated from TEMPLATE in BLUEPRINT\n\n")
	fmt.Fprintf(out, "Run\n\n")
	fmt.Fprintf(out, "  blueprint show BLUEPRINT\n\n")
	fmt.Fprintf(out, "to see possible values for TEMPLATE\n\n")
}

// Exec implements Command
func (cmd *CommandBlueprintFile) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) == 0 {
		return cmd, fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]
	args = args[1:]
	if len(args) == 0 {
		return cmd, fmt.Errorf("No file name provided")
	}
	cmd.FileName = args[0]

	args = args[1:]
	if len(args) == 0 {
		return cmd, fmt.Errorf("No template name provided")
	}
	cmd.TemplateName = args[0]

	return cmd, ctx.app.Execute(&dux.DefineBlueprintFile{
		BlueprintName: cmd.BlueprintName,
		TemplateName:  cmd.TemplateName,
		FileName:      cmd.FileName,
	})
}

// Options implements Command
func (cmd *CommandBlueprintFile) Options() *flag.FlagSet { return nil }
