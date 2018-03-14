package dux

import "path/filepath"

// ListTemplates creates a new blueprint with the given name.
type ListTemplates struct {
	BlueprintName string
}

// CommandName implements Command
func (c *ListTemplates) CommandName() string { return "list-templates" }

// ListTemplatesInFileSystem stores a blueprint in the provided blueprint store.
type ListTemplatesInFileSystem struct {
	fs     FileSystem
	events EventStore
}

// NewListTemplatesInFileSystem returns a new command handler with the given store
func NewListTemplatesInFileSystem(fs FileSystem, events EventStore) *ListTemplatesInFileSystem {
	return &ListTemplatesInFileSystem{
		fs:     fs,
		events: events,
	}
}

// Execute implements CommandHandler.
func (h *ListTemplatesInFileSystem) Execute(command Command) error {
	args := command.(*ListTemplates)
	templateDir := filepath.Join("blueprints", args.BlueprintName, "templates")
	names, err := h.fs.List(templateDir)
	if err != nil {
		return err
	}

	for _, name := range names {
		h.events.Emit(&Event{
			Name: "blueprint-template-found",
			Payload: EventPayload{
				"name":          name,
				"blueprintName": args.BlueprintName,
			},
		})
	}

	return nil
}
