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
		Describe("Inspect and edit blueprints").
		Add("template", cli.NewCommandBlueprintTemplate()).
		Add("file", cli.NewCommandBlueprintFile()).
		Add("show", cli.NewCommandBlueprintShow()).
		Add("describe", cli.NewCommandBlueprintDescribe()).
		Add("create", cli.NewCommandBlueprintCreate())

	dispatcher := cli.NewDispatchCommand(os.Args[0]).
		Add("new", cli.NewCommandNew()).
		Add("list", cli.NewCommandList()).
		Add("blueprint", blueprintCommands)

	cmd, err := cliApp.Execute(dispatcher, os.Args)
	if err != nil {
		cliApp.ShowError(err)
		if usage, ok := cmd.(cli.HasUsage); ok {
			usage.ShowUsage(cliApp.Err)
		}
	}
}
