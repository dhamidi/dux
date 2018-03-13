package dux

import "path/filepath"

// RenderBlueprint is a command for rendering a given blueprint.
type RenderBlueprint struct {
	Name        string
	Destination string
	Data        interface{}
}

// CommandName implements Command
func (cmd *RenderBlueprint) CommandName() string { return "render-blueprint" }

// RenderBlueprintToFileSystem executes a RenderBlueprint command by
// rendering the files described by the blueprint into a file system.
type RenderBlueprintToFileSystem struct {
	fs     FileSystem
	store  Store
	events EventStore
}

// NewRenderBlueprintToFileSystem returns a command handler that renders files into the provided filesystem.
func NewRenderBlueprintToFileSystem(fs FileSystem, store Store, events EventStore) *RenderBlueprintToFileSystem {
	return &RenderBlueprintToFileSystem{
		fs:     fs,
		store:  store,
		events: events,
	}
}

// Execute renders the files described by the blueprint
func (r *RenderBlueprintToFileSystem) Execute(command Command) error {
	args := command.(*RenderBlueprint)

	blueprint := new(Blueprint)
	if err := r.store.Get(args.Name, blueprint); err != nil {
		return err
	}
	templates := NewHTMLTemplateEngine(filepath.Join("blueprints", blueprint.Name, "templates"), r.fs)
	for destinationFileName, templateName := range blueprint.Files {
		var err error
		outputFilePathTemplate := filepath.Join(args.Destination, destinationFileName)
		outputFilePath, err := templates.RenderString(outputFilePathTemplate, args.Data)
		if err != nil {
			r.events.Emit(&Event{
				Name:  "render-destination-filename-failed",
				Error: err,
			})
			continue
		}
		destinationFile, err := r.fs.Create(outputFilePath)
		if err != nil {
			r.events.Emit(&Event{
				Name:  "create-destination-file-failed",
				Error: err,
			})
			continue
		}
		err = templates.RenderTemplate(destinationFile, templateName, args.Data)
		if err != nil {
			destinationFile.Close()
			r.events.Emit(&Event{
				Name:  "render-template-failed",
				Error: err,
			})
			continue
		}
		destinationFile.Close()
		r.events.Emit(&Event{
			Name: "template-rendered",
			Payload: EventPayload{
				"filename": outputFilePath,
				"template": templateName,
			},
		})
	}

	return nil
}
