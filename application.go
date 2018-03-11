package dux

import "fmt"

// Application is the entry point and context for all operations in dux.
type Application struct {
	commandHandlers map[string]CommandHandler
	FileSystem      FileSystem // the file system into which files are rendered
	Store           Store      // access to persistent storage of serialized objects.
}

// NewApplication constructs a new application instance with sensible
// defaults.
func NewApplication() *Application {
	result := &Application{
		commandHandlers: map[string]CommandHandler{},
		FileSystem:      NewInMemoryFileSystem(),
	}
	result.Store = NewFileSystemStore("blueprints", result.FileSystem)
	result.Handle("render-blueprint", NewRenderBlueprintToFileSystem(result.FileSystem, result.Store))
	result.Handle("create-blueprint", NewCreateBlueprintInFileSystem(result.Store))
	result.Handle("define-blueprint-template", NewStoreBlueprintTemplate(result.FileSystem))
	result.Handle("define-blueprint-file", NewAddFileToBlueprint(result.Store))
	return result
}

// Handle registers a command handler for the given command type.
func (app *Application) Handle(commandName string, handler CommandHandler) *Application {
	app.commandHandlers[commandName] = handler
	return app
}

// HandleFunc registers a function as a command handler for the given command.
func (app *Application) HandleFunc(commandName string, handler func(Command) error) *Application {
	app.Handle(commandName, CommandHandlerFunc(handler))
	return app
}

// Execute runs the given command in the context of this application.
func (app *Application) Execute(command Command) error {
	handler, found := app.commandHandlers[command.CommandName()]
	if !found {
		return fmt.Errorf("Command not implemented: %s", command.CommandName())
	}
	return handler.Execute(command)
}
