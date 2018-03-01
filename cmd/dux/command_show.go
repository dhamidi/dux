package main

import (
	"fmt"
	"io"

	"github.com/dhamidi/dux"
)

// CommandShow displays detailed information about a blueprint.
type CommandShow struct {
	out io.Writer
}

// CommandName implements dux.Command
func (c *CommandShow) CommandName() string { return "show" }

// CommandDescription implements dux.Command
func (c *CommandShow) CommandDescription() string { return `Display information about a blueprint` }

// Execute Display information about a blueprint
func (c *CommandShow) Execute(ctx *dux.Context, args []string) error {
	blueprintName := args[0]
	blueprint, err := ctx.LoadBlueprint(blueprintName)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.out, "Blueprint: %s\n\n", blueprintName)
	fmt.Fprintf(c.out, "%s\n\n", blueprint.Description)
	if len(blueprint.Args) > 0 {
		fmt.Fprintf(c.out, "Arguments:\n")
		for _, arg := range blueprint.Args {
			c.showArg(arg)
		}
	}
	if len(blueprint.Files) > 0 {
		fmt.Fprintf(c.out, "Generated files:\n")
		for _, file := range blueprint.Files {
			c.showFile(file)
		}
	}

	return nil
}

func (c *CommandShow) showArg(arg *dux.BlueprintArgument) {
	fmt.Fprintf(c.out, "  %s [%s]\n", arg.Name, arg.Type)
	if arg.Doc != nil {
		fmt.Fprintf(c.out, `    %s`, *arg.Doc)
	}
	fmt.Fprintf(c.out, "\n\n")
}

func (c *CommandShow) showFile(f *dux.BlueprintFileDescription) {
	fmt.Fprintf(c.out, "  %s (from template %q)\n", f.Destination, f.Template)
}
