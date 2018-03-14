package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestAddFileToBlueprint_emits_an_event_when_the_files_has_been_added(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	do(h.DefineBlueprintFile("a", "EXAMPLE", "x.tmpl"))
	h.AssertEvent(t, app.EventStore, "blueprint-file-added",
		dux.EventPayload{
			"blueprintName": "a",
			"templateName":  "x.tmpl",
			"filename":      "EXAMPLE",
		})
}
