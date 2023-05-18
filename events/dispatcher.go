package events

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrorEventNotRegistered       = errors.New("event not registered")
	ErrorHandlerNotRegistered     = errors.New("handler not registered")
	ErrorHandlerAlreadyRegistered = errors.New("handler already registered")
)

type EventDispatcher struct {
	handlers map[string][]Handler
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]Handler),
	}
}

func (d *EventDispatcher) Register(event Event, handler Handler) error {
	if handlers, ok := d.handlers[event.Name()]; ok {
		for _, h := range handlers {
			if h == handler {
				return ErrorHandlerAlreadyRegistered
			}
		}
	}

	d.handlers[event.Name()] = append(d.handlers[event.Name()], handler)
	return nil
}

func (d *EventDispatcher) Remove(event Event, handler Handler) error {
	if handlers, ok := d.handlers[event.Name()]; ok {
		for i, h := range handlers {
			if h == handler {
				d.handlers[event.Name()] = append(handlers[:i], handlers[i+1:]...)
				return nil
			}
		}

		return ErrorHandlerNotRegistered
	}

	return ErrorEventNotRegistered
}

func (d *EventDispatcher) Has(event Event, handler Handler) bool {
	if handlers, ok := d.handlers[event.Name()]; ok {
		for _, h := range handlers {
			if h == handler {
				return true
			}
		}
	}

	return false
}

func (d *EventDispatcher) Dispatch(ctx context.Context, event Event) []error {
	if handlers, ok := d.handlers[event.Name()]; ok {
		ch := make(chan error, len(handlers))
		wg := &sync.WaitGroup{}

		for _, handler := range handlers {
			h := handler
			wg.Add(1)

			go func() {
				defer wg.Done()

				err := h.Handle(ctx, event)
				if err != nil {
					ch <- err
				}
			}()
		}

		wg.Wait()
		close(ch)

		errs := make([]error, 0, len(handlers))

		for err := range ch {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			return errs
		}

		return nil
	}

	return []error{ErrorEventNotRegistered}
}

func (d *EventDispatcher) Clear() {
	d.handlers = make(map[string][]Handler)
}
