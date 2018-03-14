package cli

import (
	"flag"
	"io"
)

// Command defines the operations that need to be implemented by a CLI command.
type Command interface {
	Exec(cli *CLI, args []string) (Command, error)
	Options() *flag.FlagSet
}

// HasUsage is implemented by Commands that can describe their usage.
type HasUsage interface {
	ShowUsage(out io.Writer)
}

// HasDescription is implemented by Commands that can summarize their usage
type HasDescription interface {
	Description() string
}
