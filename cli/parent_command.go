package cli

type parentCommand struct {
	parent *DispatchCommand
}

func (cmd *parentCommand) SetParent(parent *DispatchCommand) {
	cmd.parent = parent
}

func (cmd *parentCommand) CommandPath() string {
	if cmd.parent == nil {
		return ""
	}
	return cmd.parent.CommandPath()
}
