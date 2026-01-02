package ports

import (
	"context"
	"delivery/internal/pkg/ddd"
)

type EventPublisher interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
}
