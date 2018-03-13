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
}

// TransientEventStore stores events in RAM for the runtime of the program.
type TransientEventStore struct {
	events []*Event
}

// NewTransientEventStore creates an empty transient event store.
func NewTransientEventStore() *TransientEventStore {
	return &TransientEventStore{
		events: []*Event{},
	}
}

// All returns all events that have been emitted so far.
//
// It never returns an error.
func (s *TransientEventStore) All() ([]*Event, error) {
	return s.events, nil
}

// Emit records events
//
// It never returns an error.
func (s *TransientEventStore) Emit(events ...*Event) error {
	s.events = append(s.events, events...)
	return nil
}
