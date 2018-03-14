package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestApp_CreateBlueprint_creates_a_blueprint_that_can_be_loaded_from_the_blueprint_repository(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	blueprint := new(dux.Blueprint)
	if err := app.Store.Get("a", blueprint); err != nil {
		t.Fatalf("Error loading blueprint %q: %s", "a", err)
	}
}

func TestApp_CreateBlueprint_emits_an_event_when_the_blueprint_has_been_created(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	h.AssertEvent(t, app.EventStore, "blueprint-created", dux.EventPayload{"name": "a"})
}
