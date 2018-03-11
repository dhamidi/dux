package testing

import (
	"testing"

	"github.com/dhamidi/dux"
)

func NewApp() *dux.Application {
	return dux.NewApplication()
}

func ExampleBlueprintName() string {
	return "test"
}

func RenderBlueprint(blueprintName string, context ...interface{}) *dux.RenderBlueprint {
	data := (interface{})(nil)
	if len(context) > 0 {
		data = context[0]
	}
	return &dux.RenderBlueprint{
		Name:        blueprintName,
		Destination: "staging",
		Data:        data,
	}
}

func CreateBlueprint(name string) *dux.CreateBlueprint {
	return &dux.CreateBlueprint{
		Name: name,
	}
}

func DefineBlueprintTemplate(blueprintName, templateName, contents string) *dux.DefineBlueprintTemplate {
	return &dux.DefineBlueprintTemplate{
		BlueprintName: blueprintName,
		TemplateName:  templateName,
		Contents:      contents,
	}
}

func DefineBlueprintFile(blueprintName, fileName, templateName string) *dux.DefineBlueprintFile {
	return &dux.DefineBlueprintFile{
		BlueprintName: blueprintName,
		TemplateName:  templateName,
		FileName:      fileName,
	}
}

func FailOnExecuteError(t *testing.T, h dux.CommandHandler) func(dux.Command) error {
	return func(cmd dux.Command) error {
		if err := h.Execute(cmd); err != nil {
			t.Fatalf("%s: %s", cmd.CommandName(), err)
		}
		return nil
	}
}
