package mappers

import "delivery/internal/pkg/ddd"

type ToJSONMapper interface {
	ToJSON(event ddd.DomainEvent) ([]byte, error)
}

type FromJSONMapper interface {
	FromJSON(data []byte) (ddd.DomainEvent, error)
}
