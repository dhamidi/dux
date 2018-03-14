package dux

// TransientEventStore stores events in RAM for the runtime of the program.
type TransientEventStore struct {
	events      []*Event
	subscribers []func(*Event)
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
	for _, event := range events {
		s.Notify(event)
	}
	return nil
}

// Subscribe registers a subscriber
func (s *TransientEventStore) Subscribe(subscriber func(*Event)) func() {
	position := len(s.subscribers)
	s.subscribers = append(s.subscribers, subscriber)
	return func() {
		s.subscribers[position] = func(*Event) {}
	}
}

// Notify calls all subscribers with the given event
func (s *TransientEventStore) Notify(about *Event) {
	for _, subscriber := range s.subscribers {
		subscriber(about)
	}
}
