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
	cmd := &CommandNew{BlueprintName: "command"}
	if err := cmd.Execute(app, os.Args[1:]); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
