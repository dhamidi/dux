package cli

import (
	"flag"
	"fmt"
	"io"
)

// DispatchCommand implements Command by consuming one argument and
// dispatching to subcommands.
type DispatchCommand struct {
	parent      *DispatchCommand
	name        string
	description string
	subcommands map[string]Command
}

// NewDispatchCommand creates a new dispatcher with the given name.
//
// If name is empty, no argument will be consumed.
func NewDispatchCommand(name string) *DispatchCommand {
	return &DispatchCommand{
		name:        name,
		subcommands: map[string]Command{},
	}
}

// SetParent associates this dispatcher with a parent dispatcher.
// This is mainly useful for printing the full command path when in
// ShowUsage.
//
// SetParent is called automatically by Add when adding a
// DispatchCommand.
func (cmd *DispatchCommand) SetParent(parent *DispatchCommand) {
	cmd.parent = parent
}

// Description implements HasDescription
func (cmd *DispatchCommand) Description() string { return cmd.description }

// Describe sets the description text for this dispatcher
func (cmd *DispatchCommand) Describe(desc string) *DispatchCommand {
	cmd.description = desc
	return cmd
}

// Exec implements Command by consuming the first entry in args and
// dispatching to a command with a matching name.
func (cmd *DispatchCommand) Exec(ctx *CLI, args []string) (Command, error) {
	if len(args) == 0 {
		return cmd, fmt.Errorf("No subcommand provided")
	}

	if args[0] == cmd.name {
		args = args[1:]
	}

	if len(args) == 0 {
		return cmd, fmt.Errorf("No subcommand provided")
	}

	subcommand, found := cmd.subcommands[args[0]]
	if !found {
		return cmd, fmt.Errorf("Unknown command: %q", args[0])
	}

	return ctx.Execute(subcommand, args[1:])
}

// ShowUsage implements HasUsage by listing all subcommands
func (cmd *DispatchCommand) ShowUsage(out io.Writer) {
	if len(cmd.subcommands) == 0 {
		return
	}

	fmt.Fprintf(out, "\nUsage: %s SUBCOMMAND\n\n", cmd.CommandPath())

	if desc := cmd.Description(); desc != "" {
		fmt.Fprintf(out, "%s\n\n", desc)
	}

	fmt.Fprintf(out, "Available commands:\n")
	longestSubcommandName := ""
	for name := range cmd.subcommands {
		if len(name) > len(longestSubcommandName) {
			longestSubcommandName = name
		}
	}

	subcommandFormat := fmt.Sprintf("  %%-%ds", len(longestSubcommandName))
	for name, command := range cmd.subcommands {
		fmt.Fprintf(out, subcommandFormat, name)
		if description, ok := command.(HasDescription); ok {
			fmt.Fprintf(out, "  %s", description.Description())
		}
		fmt.Fprintf(out, "\n")
	}
}

// Options implements Command by returning the options for the given subcommand
func (cmd *DispatchCommand) Options() *flag.FlagSet { return nil }

// 	if len(args) == 0 {
// 		return nil
// 	}

// 	if args[0] == cmd.name {
// 		args = args[1:]
// 	}

// 	if len(args) == 0 {
// 		return nil
// 	}

// 	subcommand, found := cmd.subcommands[args[0]]
// 	if !found {
// 		return nil
// 	}

// 	return subcommand.Options(args[1:])
// }

// Add defines a new subcommand
func (cmd *DispatchCommand) Add(name string, command Command) *DispatchCommand {
	if child, isChild := command.(interface {
		SetParent(cmd *DispatchCommand)
	}); isChild {
		child.SetParent(cmd)
	}
	cmd.subcommands[name] = command
	return cmd
}

// CommandPath returns the full path to this subcommand as a list of
// space separated words.
func (cmd *DispatchCommand) CommandPath() string {
	if cmd.parent != nil {
		return cmd.parent.CommandPath() + " " + cmd.name
	}

	return cmd.name
}
