package dux_test

import (
	"testing"

	"github.com/dhamidi/dux"
)

func TestTransientEventStore_Emit_notifies_subscribers(t *testing.T) {
	eventStore := dux.NewTransientEventStore()
	notified := false
	eventStore.Subscribe(func(event *dux.Event) {
		notified = true
	})

	eventStore.Emit(&dux.Event{Name: "test"})

	if !notified {
		t.Fatal("Subscriber not notified")
	}
}
