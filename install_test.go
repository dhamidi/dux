package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestInstallInFileSystem_emits_an_event_for_each_renamed_file(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	do(h.DefineBlueprintFile("a", "EXAMPLE", "x.tmpl"))
	do(h.RenderBlueprint("a", map[string]interface{}{"n": 1}))
	do(h.Install("staging/EXAMPLE", "./example"))
	h.AssertEvent(t, app.EventStore, "file-renamed",
		dux.EventPayload{
			"to":   "./example",
			"from": "staging/EXAMPLE",
		})

}

func TestInstallInFileSystem_emits_an_event_for_file_it_failed_to_rename(t *testing.T) {
	app := h.NewApp()
	failingFS := h.NewFailingFileSystem(app.FileSystem)
	app.FileSystem = failingFS
	app.Init()
	failingFS.Fail("rename", "staging/EXAMPLE")
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DefineBlueprintTemplate("a", "x.tmpl", "{{.n}}"))
	do(h.DefineBlueprintFile("a", "EXAMPLE", "x.tmpl"))
	do(h.RenderBlueprint("a", map[string]interface{}{"n": 1}))
	do(h.Install("staging/EXAMPLE", "./example"))
	h.AssertEvent(t, app.EventStore, "file-rename-failed",
		dux.EventPayload{
			"to":   "./example",
			"from": "staging/EXAMPLE",
		})
}
