package dux

// DescribeBlueprint defines a template that should be associated with the blueprint.
type DescribeBlueprint struct {
	BlueprintName string
	Description   string
}

// CommandName implements Command
func (c *DescribeBlueprint) CommandName() string { return "describe-blueprint" }

// SetBlueprintDescription loads the blueprint from the store, sets the description of the blueprint and stores the blueprint again.
type SetBlueprintDescription struct {
	store  Store
	events EventStore
}

// NewSetBlueprintDescription returns a new command handler with the given file system.
func NewSetBlueprintDescription(store Store, events EventStore) *SetBlueprintDescription {
	return &SetBlueprintDescription{store: store, events: events}
}

// Execute implements CommandHandler
func (h *SetBlueprintDescription) Execute(command Command) error {
	args := command.(*DescribeBlueprint)
	blueprint := new(Blueprint)
	if err := h.store.Get(args.BlueprintName, blueprint); err != nil {
		return err
	}
	blueprint.SetDescription(args.Description)
	err := h.store.Put(args.BlueprintName, blueprint)
	if err == nil {
		h.events.Emit(&Event{
			Name: "blueprint-description-set",
			Payload: EventPayload{
				"blueprintName": args.BlueprintName,
				"description":   args.Description,
			},
		})
	}
	return err
}
