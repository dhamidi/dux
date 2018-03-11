package testing

import "github.com/dhamidi/dux"

func NewApp() *dux.Application {
	return dux.NewApplication()
}

func RenderBlueprint() *dux.RenderBlueprint {
	return &dux.RenderBlueprint{
		Destination: "staging",
	}
}

func ExampleBlueprintFilenames() []string {
	return []string{"test-file"}
}
