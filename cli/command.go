package cli

import "flag"

// Command defines the operations that need to be implemented by a CLI command.
type Command interface {
	Exec(cli *CLI, args []string) error
	Options() *flag.FlagSet
}
