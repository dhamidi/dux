package testing

import (
	"fmt"
	"reflect"
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

func DescribeBlueprint(name string, description string) *dux.DescribeBlueprint {
	return &dux.DescribeBlueprint{
		BlueprintName: name,
		Description:   description,
	}
}

func Install(pairs ...string) *dux.Install {
	sources := []string{}
	destinations := []string{}
	if len(pairs)%2 != 0 {
		panic(fmt.Sprintf("Install: uneven number of arguments in %#v", pairs))
	}
	for i, f := range pairs {
		if i%2 == 0 {
			sources = append(sources, f)
		} else {
			destinations = append(destinations, f)
		}
	}

	return &dux.Install{
		Sources:      sources,
		Destinations: destinations,
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

func AssertEvent(t *testing.T, events dux.EventStore, eventName string, expectedPayload dux.EventPayload) {
	t.Helper()
	eventNames := []string{}
	allEvents, err := events.All()
	if err != nil {
		t.Fatalf("AssertEvent: failed to fetch events from event store: %s", err)
	}
	for _, event := range allEvents {
		eventNames = append(eventNames, event.Name)
		if event.Name != eventName {
			continue
		}

		for key, expectedValue := range expectedPayload {
			actualValue := event.Payload[key]
			if reflect.DeepEqual(expectedValue, actualValue) {
				continue
			}
			t.Fatalf("event payload mismatch: want %#v, got %#v",
				expectedValue,
				actualValue,
			)
		}
		return
	}

	t.Fatalf("Event %q not found in %v", eventName, eventNames)
}
