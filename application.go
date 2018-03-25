package dux

import "fmt"

// Application is the entry point and context for all operations in dux.
type Application struct {
	commandHandlers map[string]CommandHandler
	FileSystem      FileSystem // the file system into which files are rendered
	Store           Store      // access to persistent storage of serialized objects.
	EventStore      EventStore // access to events that have been emitted by commands
}

// NewApplication constructs a new application instance with sensible
// defaults.
func NewApplication() *Application {
	result := &Application{
		commandHandlers: map[string]CommandHandler{},
		FileSystem:      NewInMemoryFileSystem(),
		EventStore:      NewTransientEventStore(),
	}
	return result.Init()
}

// Init installs all default command handlers after clearing all registered command handlers.
func (app *Application) Init() *Application {
	app.commandHandlers = map[string]CommandHandler{}
	app.Store = NewFileSystemStore("blueprints", app.FileSystem)
	app.Handle("render-blueprint", NewRenderBlueprintToFileSystem(app.FileSystem, app.Store, app.EventStore))
	app.Handle("create-blueprint", NewCreateBlueprintInFileSystem(app.Store, app.EventStore))
	app.Handle("define-blueprint-template", NewStoreBlueprintTemplate(app.FileSystem, app.EventStore))
	app.Handle("define-blueprint-file", NewAddFileToBlueprint(app.Store, app.EventStore))
	app.Handle("describe-blueprint", NewSetBlueprintDescription(app.Store, app.EventStore))
	app.Handle("list-templates", NewListTemplatesInFileSystem(app.FileSystem, app.EventStore))
	app.Handle("install", NewInstallInFileSystem(app.FileSystem, app.EventStore))
	return app
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
