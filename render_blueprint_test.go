package dux_test

import (
	"testing"

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
