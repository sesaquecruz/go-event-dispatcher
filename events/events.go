// Package events contains interfaces and implementations to deal with events.
package events

import "context"

// Event name type
type EventName string

// Event interface.
type Event interface {
	Name() EventName
	Payload() any
}

// Event Handler interface.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Event Dispatcher interface.
type Dispatcher interface {
	Register(event EventName, handler Handler) error
	Remove(event EventName, handler Handler) error
	Has(event EventName, handler Handler) bool
	Dispatch(ctx context.Context, event Event) error
	Clear()
}
