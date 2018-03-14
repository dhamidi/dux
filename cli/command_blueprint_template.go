package cli

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/dhamidi/dux"
)

// CommandBlueprintTemplate is a CLI command for rendering a blueprint.
type CommandBlueprintTemplate struct {
	BlueprintName string
	TemplateName  string
	Contents      string
}

// NewCommandBlueprintTemplate creates a new, empty instance of this command.
func NewCommandBlueprintTemplate() *CommandBlueprintTemplate {
	return &CommandBlueprintTemplate{}
}

// Exec implements Command
func (cmd *CommandBlueprintTemplate) Exec(ctx *CLI, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]
	args = args[1:]
	if len(args) == 0 {
		return fmt.Errorf("No template name provided")
	}
	cmd.TemplateName = args[0]

	if err := cmd.ReadContents(ctx); err != nil {
		return err
	}

	return ctx.app.Execute(&dux.DefineBlueprintTemplate{
		BlueprintName: cmd.BlueprintName,
		TemplateName:  cmd.TemplateName,
		Contents:      cmd.Contents,
	})
}

// Options implements Command
func (cmd *CommandBlueprintTemplate) Options() *flag.FlagSet {
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.StringVar(&cmd.Contents, "contents", "", "Template contents")
	return flags
}

// ReadContents reads contents from the template file from the input associated with the provided context.
func (cmd *CommandBlueprintTemplate) ReadContents(ctx *CLI) error {
	if cmd.Contents != "" {
		return nil
	}

	contents, err := ioutil.ReadAll(ctx.in)
	if err != nil {
		return err
	}

	cmd.Contents = string(contents)
	return nil
}
