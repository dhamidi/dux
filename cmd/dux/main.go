package main

import (
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
	cliApp.Execute(cli.NewCommandNew(), os.Args[1:])
}
