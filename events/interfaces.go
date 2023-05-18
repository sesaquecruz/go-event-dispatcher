package events

import "context"

type Event interface {
	Name() string
	Payload() any
}

type Handler interface {
	Handle(ctx context.Context, event Event) error
}

type Dispatcher interface {
	Register(event Event, handler Handler) error
	Remove(event Event, handler Handler) error
	Has(event Event, handler Handler) bool
	Dispatch(ctx context.Context, event Event) error
	Clear()
}
