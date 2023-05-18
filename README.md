# Go Event Dispatcher

This project contains a Go lib to assist dealing with events. It provides a set of interfaces for events, handlers, and dispatchers, as well as an implementation of an event dispatcher that supports registering events and their handlers. Event dispatchers are useful in systems that commonly use asynchronous communications, for example, in event-driven systems that utilize message brokers such as RabbitMQ and Kafka to communication between them.

## Installation

To install, you can use the go get command:

```
go get github.com/sesaquecruz/go-event-dispatcher
```

## Usage

1. Import the lib in your Go code:

```
import "github.com/sesaquecruz/go-event-dispatcher/events"
```

2. Define your events following the `Event` interface.

```
type Event interface {
	Name() EventName
	Payload() any
}
```

```
// An event example
type InvoiceEvent struct {
	name    events.EventName
	payload InvoicePayload
}

func (e *InvoiceEvent) Name() events.EventName {
	return e.name
}

func (e *InvoiceEvent) Payload() any {
	return e.payload
}
```

```
// An event example
type OrderEvent struct {
	name    events.EventName
	payload OrderPayload
}

func (e *OrderEvent) Name() events.EventName {
	return e.name
}

func (e *OrderEvent) Payload() any {
	return e.payload
}
```

3. Define your handlers following the `Handler` interface.

```
type Handler interface {
	Handle(ctx context.Context, event Event) error
}
```

```
// A handler example
type SaleHandler struct {
	// Communication channel with sale system
	// ...
}

func (h *SaleHandler) Handle(ctx context.Context, event events.Event) error {
	// Sending message implementation
	// ...
}
```

```
// A handler example
type DeliveryHandler struct {
	// Communication channel with delivery system
	// ...
}

func (h *DeliveryHandler) Handle(ctx context.Context, event events.Event) error {
	// Sending message implementation
	// ...
}
```

4. Create an event dispatch, register the events and its handlers, and dispatch them when necessary.

```
func main() {
	// Create event names (types), events, handlers, and contexts
	// ...
	
	dispatcher := events.NewEventDispatcher()

	dispatcher.Register(invoiceEventType, saleHandler)

	dispatcher.Register(orderEventType, saleHandler)
	dispatcher.Register(orderEventType, deliveryHandler)

	dispatcher.Dispatch(ctx, invoiceEvent)
	dispatcher.Dispatch(ctx, orderEvent)
}
```

## Contributing

Contributions to this project are welcome. If you encounter any issues or have ideas for enhancements, feel free to open an issue or submit a pull request.

## License
This project is licensed under the MIT License. Please see the [LICENSE](./LICENSE) file for more details.
