// Package events contains interfaces and implementations to deal with events.
package events

import "context"

// Event interface.
type Event interface {
	Name() string
	Payload() any
}

// Event Handler interface.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Event Dispatcher interface.
type Dispatcher interface {
	Register(event Event, handler Handler) error
	Remove(event Event, handler Handler) error
	Has(event Event, handler Handler) bool
	Dispatch(ctx context.Context, event Event) error
	Clear()
}
