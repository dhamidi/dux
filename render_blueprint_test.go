package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestApp_RenderBlueprint_renders_blueprints_that_have_been_previously_created(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "Test"))
	do(h.DefineBlueprintFile("a", "x-file", "x.tmpl"))
	do(h.RenderBlueprint("a"))
	h.AssertFileContents(t, app.FileSystem, "staging/x-file", "Test")
}

func TestApp_RenderBlueprint_provides_the_given_context_to_the_template_object(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	do(h.DefineBlueprintFile("a", "x-file", "x.tmpl"))
	do(h.RenderBlueprint("a", map[string]interface{}{"n": 1}))
	h.AssertFileContents(t, app.FileSystem, "staging/x-file", "1")
}

func TestApp_RenderBlueprint_renders_destination_file_names_as_templates(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	do(h.DefineBlueprintFile("a", "{{.n}}-file", "x.tmpl"))
	do(h.RenderBlueprint("a", map[string]interface{}{"n": 1}))
	h.AssertFileContents(t, app.FileSystem, "staging/1-file", "1")
}

func TestApp_RenderBlueprint_emits_an_event_for_each_successfully_rendered_file(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	do(h.DefineBlueprintFile("a", "{{.n}}-file", "x.tmpl"))
	do(h.RenderBlueprint("a", map[string]interface{}{"n": 1}))
	h.AssertEvent(t, app.EventStore, "template-rendered",
		dux.EventPayload{"filename": "staging/1-file"},
	)
}
