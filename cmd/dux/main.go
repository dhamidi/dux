package main

import (
	"fmt"
	"os"

	"github.com/dhamidi/dux"
)

type CommandNew struct {
	BlueprintName string
}

// CommandName implements dux.Command
func (c *CommandNew) CommandName() string { return "new" }

// CommandDescription implements dux.Command
func (c *CommandNew) CommandDescription() string { return `Generate files from a blueprint.` }

// Execute renders templates from the blueprint
func (c *CommandNew) Execute(ctx *dux.Context, args []string) error {
	blueprint, err := ctx.LoadBlueprint(c.BlueprintName)
	if err != nil {
		return err
	}
	if err := blueprint.ParseArgs(args); err != nil {
		return err
	}
	if err := ctx.GatherData("dux.json"); err != nil {
		return err
	}

	return blueprint.Render(ctx)
}

func main() {
	dux := dux.NewContextFromEnvironment(dux.SystemEnvironment)
	if err := dux.GatherData("dux.json"); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	cmd := &CommandNew{BlueprintName: "command"}
	if err := cmd.Execute(dux, os.Args[1:]); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
