package ports

import (
	"context"
	"delivery/internal/pkg/ddd"
)

type EventProducer[TDomainEvent ddd.DomainEvent] interface {
	Publish(ctx context.Context, domainEvent TDomainEvent) error
	Close() error
}
