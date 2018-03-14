package dux

// EventPayload is an unordered set of key-value pairs that carries event-specific information
type EventPayload map[string]interface{}

// Event documents an action that has taken place.
type Event struct {
	Name    string
	Payload EventPayload
	Error   error
}

// EventStore provides access to events that have been emitted during the execution of commands.
type EventStore interface {
	// All returns all events that have been previously emitted
	All() ([]*Event, error)
	// Emit stores events in the event store and notifies
	// subscribers about new events
	Emit(events ...*Event) error
	// Subscribe registers a function that is will be called for
	// every event that has been successfully emitted.
	Subscribe(func(*Event))
}
