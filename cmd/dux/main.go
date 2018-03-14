package main

import (
	"fmt"
	"os"

	"github.com/dhamidi/dux"
	"github.com/dhamidi/dux/cli"
)

func main() {
	app := dux.NewApplication()
	app.Execute(&dux.CreateBlueprint{Name: "example"})
	app.Execute(&dux.DefineBlueprintTemplate{BlueprintName: "example", TemplateName: "example.tmpl", Contents: "hello, world"})
	app.Execute(&dux.DefineBlueprintFile{BlueprintName: "example", FileName: "EXAMPLE", TemplateName: "example.tmpl"})

	cliApp := cli.NewCLI(app)
	app.EventStore.Subscribe(func(e *dux.Event) {
		if e.Error == nil {
			fmt.Printf("%s %v\n", e.Name, e.Payload)
			return
		}
		cliApp.ShowError(e.Error)
	})
	cliApp.Execute(cli.NewCommandNew(), os.Args[1:])
}
