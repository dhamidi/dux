package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestStoreBlueprintTemplate_emits_an_event_when_the_template_has_been_rendered(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	h.AssertEvent(t, app.EventStore, "blueprint-template-defined",
		dux.EventPayload{
			"blueprintName": "a",
			"templateName":  "x.tmpl",
		})
}
