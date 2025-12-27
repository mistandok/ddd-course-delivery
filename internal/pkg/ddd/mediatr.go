package ddd

import (
	"context"

	"github.com/mehdihadeli/go-mediatr"
)

type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

type eventPublisher struct {
}

func NewEventPublisher() EventPublisher {
	return &eventPublisher{}
}

func (e *eventPublisher) Publish(ctx context.Context, event DomainEvent) error {
	return mediatr.Publish(ctx, event)
}
