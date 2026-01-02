package ddd

import (
	"context"
	modelEvent "delivery/internal/core/domain/model/event"
	"delivery/internal/pkg/ddd"

	"github.com/mehdihadeli/go-mediatr"
)

type EventPublisher interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
}

type eventPublisher struct {
}

func NewEventPublisher() EventPublisher {
	return &eventPublisher{}
}

func (e *eventPublisher) Publish(ctx context.Context, domainEvent ddd.DomainEvent) error {
	switch domainEvent := domainEvent.(type) {
	case *modelEvent.OrderCreated:
		err := mediatr.Publish(ctx, domainEvent)
		if err != nil {
			return err
		}
	case *modelEvent.OrderCompleted:
		err := mediatr.Publish(ctx, domainEvent)
		if err != nil {
			return err
		}
	}

	return nil
}
