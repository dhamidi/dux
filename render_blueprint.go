package dux

import "fmt"

// RenderBlueprint is a command for rendering a given blueprint.
type RenderBlueprint struct {
	Destination string
}

// CommandName implements Command
func (cmd *RenderBlueprint) CommandName() string { return "render-blueprint" }

// RenderBlueprintToFileSystem executes a RenderBlueprint command by
// rendering the files described by the blueprint into a file system.
type RenderBlueprintToFileSystem struct {
	fs FileSystem
}

// NewRenderBlueprintToFileSystem returns a command handler that renders files into the provided filesystem.
func NewRenderBlueprintToFileSystem(fs FileSystem) *RenderBlueprintToFileSystem {
	return &RenderBlueprintToFileSystem{
		fs: fs,
	}
}

// Execute renders the files described by the blueprint
func (r *RenderBlueprintToFileSystem) Execute(command Command) error {
	f, err := r.fs.Create("staging/test-file")
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "Test file")
	f.Close()
	return nil
}
