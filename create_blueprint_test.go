package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestApp_CreateBlueprint_creates_a_blueprint_that_can_be_loaded_from_the_blueprint_repository(t *testing.T) {
	app := h.NewApp()
	cmd := h.CreateBlueprint("a")
	if err := app.Execute(cmd); err != nil {
		t.Fatalf("%s: %s", cmd.CommandName(), err)
	}

	blueprint := new(dux.Blueprint)
	if err := app.Store.Get("a", blueprint); err != nil {
		t.Fatalf("Error loading blueprint %q: %s", "a", err)
	}
}
