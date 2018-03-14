package dux

// Install moves files from sources to destinations and emits events
// about the progress.
type Install struct {
	Sources      []string
	Destinations []string
}

// CommandName implements Command
func (c *Install) CommandName() string { return "install" }

// InstallInFileSystem moves files in the given file system and emits
// events about the progress.
type InstallInFileSystem struct {
	fs     FileSystem
	events EventStore
}

// NewInstallInFileSystem returns a new command handler
func NewInstallInFileSystem(fs FileSystem, events EventStore) *InstallInFileSystem {
	return &InstallInFileSystem{
		fs:     fs,
		events: events,
	}
}

// Execute implements CommandHandler.
func (h *InstallInFileSystem) Execute(command Command) error {
	args := command.(*Install)
	for i, source := range args.Sources {
		err := h.fs.Rename(source, args.Destinations[i])
		if err == nil {
			h.events.Emit(&Event{
				Name: "file-renamed",
				Payload: EventPayload{
					"from": source,
					"to":   args.Destinations[i],
				},
			})
		} else {
			h.events.Emit(&Event{
				Name:  "file-rename-failed",
				Error: err,
				Payload: EventPayload{
					"from": source,
					"to":   args.Destinations[i],
				},
			})
		}
	}

	return nil
}
