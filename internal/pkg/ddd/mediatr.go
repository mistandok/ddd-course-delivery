package ddd

import "context"

type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

type Mediatr interface {
	Subscribe(handler EventHandler, events ...DomainEvent)
	Publish(ctx context.Context, event DomainEvent) error
}

type mediatr struct {
	handlers map[string][]EventHandler
}

func NewMediatr() Mediatr {
	return &mediatr{handlers: make(map[string][]EventHandler)}
}

func (e *mediatr) Subscribe(handler EventHandler, events ...DomainEvent) {
	for _, event := range events {
		handlers := e.handlers[event.GetName()]
		handlers = append(handlers, handler)
		e.handlers[event.GetName()] = handlers
	}
}

func (e *mediatr) Publish(ctx context.Context, event DomainEvent) error {
	for _, handler := range e.handlers[event.GetName()] {
		err := handler.Handle(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}
