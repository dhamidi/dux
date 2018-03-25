package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dhamidi/dux"
)

// CommandBlueprintTemplate is a CLI command for rendering a blueprint.
type CommandBlueprintTemplate struct {
	*parentCommand

	BlueprintName string
	TemplateName  string
	Contents      string
}

// NewCommandBlueprintTemplate creates a new, empty instance of this command.
func NewCommandBlueprintTemplate() *CommandBlueprintTemplate {
	return &CommandBlueprintTemplate{
		parentCommand: new(parentCommand),
	}
}

// Exec implements Command
func (cmd *CommandBlueprintTemplate) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) == 0 {
		return cmd, fmt.Errorf("No blueprint name provided")
	}
	cmd.BlueprintName = args[0]
	args = args[1:]
	if len(args) == 0 {
		return cmd, fmt.Errorf("No template name provided")
	}
	cmd.TemplateName = args[0]

	if err := cmd.ReadContents(ctx); err != nil {
		return cmd, err
	}

	return cmd, ctx.app.Execute(&dux.DefineBlueprintTemplate{
		BlueprintName: cmd.BlueprintName,
		TemplateName:  cmd.TemplateName,
		Contents:      cmd.Contents,
	})
}

// Description implements HasDescription
func (cmd *CommandBlueprintTemplate) Description() string { return `Define a template for a blueprint` }

// ShowUsage implements HasUsage
func (cmd *CommandBlueprintTemplate) ShowUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s template [--contents='...'] BLUEPRINT TEMPLATE-NAME\n\n", cmd.CommandPath())
	fmt.Fprintf(out, "Adds a template called TEMPLATE-NAME to BLUEPRINT.\n\n")
	fmt.Fprintf(out, "If no template content is provided via the contents option, the template content is read from stdin.\n\n")
	fmt.Fprintf(out, "Options:\n")
	fmt.Fprintf(out, "  --contents='...'   Set template content")
	fmt.Fprintf(out, "\n")
}

// Options implements Command
func (cmd *CommandBlueprintTemplate) Options() *flag.FlagSet {
	flags := flag.NewFlagSet("blueprint template", flag.ContinueOnError)
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
