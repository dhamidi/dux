package main

import (
	"os"

	"github.com/dhamidi/dux"
	"github.com/dhamidi/dux/cli"
)

func main() {
	app := dux.NewApplication()
	cliApp := cli.NewCLI(app)
	cliApp.Execute(cli.NewCommandNew(), os.Args[1:])
}
