package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhamidi/dux"
)

// CommandNewBlueprint sets up the directory structure and files for a new blueprint.
type CommandNewBlueprint struct {
	flags       *flag.FlagSet
	name        string
	description string
}

// NewCommandNewBlueprint returns a new command instance with all available flags accepted by the command preconfigured.
func NewCommandNewBlueprint() *CommandNewBlueprint {
	result := &CommandNewBlueprint{}
	result.flags = flag.NewFlagSet(result.CommandName(), flag.ContinueOnError)
	result.flags.StringVar(&result.name, "name", "", "Name of the blueprint")
	result.flags.StringVar(&result.description, "description", "", "Description of the blueprint")
	return result
}

// CommandName implements dux.Command
func (c *CommandNewBlueprint) CommandName() string { return "new-blueprint" }

// CommandDescription implements dux.Command
func (c *CommandNewBlueprint) CommandDescription() string { return `Initialize a new blueprint` }

// Execute Initialize a new blueprint
func (c *CommandNewBlueprint) Execute(ctx *dux.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return err
	}
	if c.name == "" {
		return fmt.Errorf("missing argument: name")
	}
	blueprint := &dux.Blueprint{
		Name:        c.name,
		Description: c.description,
		Args:        []*dux.BlueprintArgument{},
		Files:       []*dux.BlueprintFileDescription{},
	}
	manifestJSON, err := json.MarshalIndent(blueprint, "", "  ")
	if err != nil {
		return nil
	}

	blueprintDir := filepath.Join(ctx.BaseDir, "blueprints", blueprint.Name)
	manifestPath := filepath.Join(blueprintDir, "manifest.json")
	if err := os.MkdirAll(blueprintDir, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	if err := ioutil.WriteFile(manifestPath, manifestJSON, 0644); err != nil {
		return err
	}
	ctx.Log("", "create", strings.TrimPrefix(manifestPath, ctx.BaseDir+"/"))
	if err := os.Mkdir(filepath.Join(blueprintDir, "templates"), 0755); err != nil {
		return err
	}
	ctx.Log("", "create", strings.TrimPrefix(filepath.Join(blueprintDir, "templates"), ctx.BaseDir+"/"))
	return nil
}
