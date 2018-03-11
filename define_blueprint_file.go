package dux

// DefineBlueprintFile defines a template that should be associated with the blueprint.
type DefineBlueprintFile struct {
	BlueprintName string
	FileName      string
	TemplateName  string
}

// CommandName implements Command
func (c *DefineBlueprintFile) CommandName() string { return "define-blueprint-file" }

// AddFileToBlueprint loads the blueprint from the store, adds the given file and then stores the blueprint again.
type AddFileToBlueprint struct {
	store Store
}

// NewAddFileToBlueprint returns a new command handler with the given file system.
func NewAddFileToBlueprint(store Store) *AddFileToBlueprint {
	return &AddFileToBlueprint{store: store}
}

// Execute implements CommandHandler
func (h *AddFileToBlueprint) Execute(command Command) error {
	args := command.(*DefineBlueprintFile)
	blueprint := new(Blueprint)
	if err := h.store.Get(args.BlueprintName, blueprint); err != nil {
		return err
	}
	blueprint.DefineFile(args.FileName, args.TemplateName)
	return h.store.Put(args.BlueprintName, blueprint)
}
