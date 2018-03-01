package main

import (
	"fmt"
	"os"

	"github.com/dhamidi/dux"
)

func main() {
	app := dux.NewContextFromEnvironment(dux.SystemEnvironment)
	if err := app.GatherData("dux.json"); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	app.App.DefineCommand(new(CommandNew))
	app.App.DefineCommand(&CommandShow{out: os.Stdout})
	app.App.DefineCommand(&CommandList{out: os.Stdout})
	if len(os.Args) == 1 {
		app.App.Usage(os.Stdout)
		return
	}
	app.App.RunCommand(os.Args[1], os.Args[2:])
}
