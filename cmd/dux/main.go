package main

import (
	"fmt"
	"os"

	"github.com/dhamidi/dux"
	"github.com/dhamidi/dux/cli"
)

func main() {
	app := dux.NewApplication()
	app.FileSystem = dux.NewOnDiskFileSystem()
	app.Init()
	cliApp := cli.NewCLI(app)
	app.EventStore.Subscribe(func(e *dux.Event) {
		if e.Name == "blueprint-template-found" {
			return
		}
		if e.Error == nil {
			fmt.Printf("%s %v\n", e.Name, e.Payload)
			return
		}
		cliApp.ShowError(e.Error)
	})

	blueprintCommands := cli.NewDispatchCommand("blueprint").
		Add("template", cli.NewCommandBlueprintTemplate()).
		Add("file", cli.NewCommandBlueprintFile()).
		Add("show", cli.NewCommandBlueprintShow())

	dispatcher := cli.NewDispatchCommand("").
		Add("new", cli.NewCommandNew()).
		Add("blueprint", blueprintCommands)

	cliApp.Execute(dispatcher, os.Args[1:])
}
