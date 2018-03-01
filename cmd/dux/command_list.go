package main

import (
	"fmt"
	"io"

	"github.com/dhamidi/dux"
)

// CommandList writes the names of all available blueprints to the given io.Writer
type CommandList struct {
	out io.Writer
}

// CommandName implements dux.Command
func (c *CommandList) CommandName() string { return "list" }

// CommandDescription implements dux.Command
func (c *CommandList) CommandDescription() string { return `List all available blueprints` }

// Execute List all available blueprints
func (c *CommandList) Execute(ctx *dux.Context, args []string) error {
	blueprintNames, err := ctx.ListBlueprints()
	if err != nil {
		return err
	}
	for _, name := range blueprintNames {
		fmt.Fprintf(c.out, "%s\n", name)
	}
	return nil
}
