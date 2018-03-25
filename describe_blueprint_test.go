package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
	h "github.com/dhamidi/dux/testing"
)

func TestDescribeBlueprint_emits_an_event_when_the_blueprint_has_been_saved(t *testing.T) {
	app := h.NewApp()
	do := h.FailOnExecuteError(t, app)
	do(h.CreateBlueprint("a"))
	do(h.DescribeBlueprint("a", "A test"))
	h.AssertEvent(t, app.EventStore, "blueprint-description-set",
		dux.EventPayload{
			"blueprintName": "a",
			"description":   "A test",
		})
}
