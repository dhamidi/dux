package dux

import (
	"io"
	"path/filepath"
)

// DefineBlueprintTemplate defines a template that should be associated with the blueprint.
type DefineBlueprintTemplate struct {
	BlueprintName string
	TemplateName  string
	Contents      string
}

// CommandName implements Command
func (c *DefineBlueprintTemplate) CommandName() string { return "define-blueprint-template" }

// StoreBlueprintTemplate writes the given template into a subdirectory of the given blueprint's directory.
type StoreBlueprintTemplate struct {
	fs     FileSystem
	events EventStore
}

// NewStoreBlueprintTemplate returns a new command handler with the given file system.
func NewStoreBlueprintTemplate(fs FileSystem, events EventStore) *StoreBlueprintTemplate {
	return &StoreBlueprintTemplate{
		fs:     fs,
		events: events,
	}
}

// Execute implements CommandHandler
func (h *StoreBlueprintTemplate) Execute(command Command) error {
	args := command.(*DefineBlueprintTemplate)
	destinationFile := filepath.Join("blueprints", args.BlueprintName, "templates", args.TemplateName)
	out, err := h.fs.Create(destinationFile)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.WriteString(out, args.Contents)
	if err == nil {
		h.events.Emit(&Event{
			Name: "blueprint-template-defined",
			Payload: EventPayload{
				"blueprintName": args.BlueprintName,
				"templateName":  args.TemplateName,
			},
		})
	}
	return err
}
