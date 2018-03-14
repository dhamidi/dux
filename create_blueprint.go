package dux

// CreateBlueprint creates a new blueprint with the given name.
type CreateBlueprint struct {
	Name string
}

// CommandName implements Command
func (c *CreateBlueprint) CommandName() string { return "create-blueprint" }

// CreateBlueprintInFileSystem stores a blueprint in the provided blueprint store.
type CreateBlueprintInFileSystem struct {
	store  Store
	events EventStore
}

// NewCreateBlueprintInFileSystem returns a new command handler with the given store
func NewCreateBlueprintInFileSystem(store Store, events EventStore) *CreateBlueprintInFileSystem {
	return &CreateBlueprintInFileSystem{
		store:  store,
		events: events,
	}
}

// Execute implements CommandHandler.
func (h *CreateBlueprintInFileSystem) Execute(command Command) error {
	createBlueprint := command.(*CreateBlueprint)
	blueprint := &Blueprint{
		Name: createBlueprint.Name,
	}
	err := h.store.Put(blueprint.Name, blueprint)
	if err == nil {
		h.events.Emit(&Event{
			Name:    "blueprint-created",
			Payload: EventPayload{"name": createBlueprint.Name},
		})
	}
	return err
}
