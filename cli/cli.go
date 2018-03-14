package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/dhamidi/dux"
)

// CLI models the dux CLI application and serves as a container for
// commands and associated command handlers.
type CLI struct {
	app   *dux.Application
	in    io.Reader
	out   io.Writer
	Err   io.Writer
	depth int
}

// NewCLI creates a new CLI application wrapping the provided dux
// instance and connected to os.Stdout, os.Stderr and os.Stdin by
// default.
func NewCLI(app *dux.Application) *CLI {
	return &CLI{
		app:   app,
		in:    os.Stdin,
		out:   os.Stdout,
		Err:   os.Stderr,
		depth: 0,
	}
}

// Execute runs a given command with the given arguments.  Any errors returned by the command are shown to the user
func (cli *CLI) Execute(cmd Command, args []string) (Command, error) {
	options := cmd.Options()
	if options != nil {
		options.Parse(args)
		args = options.Args()
	}
	return cmd.Exec(cli, args)
}

// ShowError displays an error to the user
func (cli *CLI) ShowError(err error) {
	fmt.Fprintf(cli.Err, "Error: %s\n", err)
}
