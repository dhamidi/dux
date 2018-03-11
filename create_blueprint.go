package dux

// CreateBlueprint creates a new blueprint with the given name.
type CreateBlueprint struct {
	Name string
}

// CommandName implements Command
func (c *CreateBlueprint) CommandName() string { return "create-blueprint" }

// CreateBlueprintInFileSystem stores a blueprint in the provided blueprint store.
type CreateBlueprintInFileSystem struct {
	store Store
}

// NewCreateBlueprintInFileSystem returns a new command handler with the given store
func NewCreateBlueprintInFileSystem(store Store) *CreateBlueprintInFileSystem {
	return &CreateBlueprintInFileSystem{
		store: store,
	}
}

// Execute implements CommandHandler.
func (h *CreateBlueprintInFileSystem) Execute(command Command) error {
	createBlueprint := command.(*CreateBlueprint)
	blueprint := &Blueprint{
		Name: createBlueprint.Name,
	}
	return h.store.Put(blueprint.Name, blueprint)
}
