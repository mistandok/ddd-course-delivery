package ddd

import (
	"github.com/google/uuid"
)

type DomainEvent interface {
	GetID() uuid.UUID
	GetName() string
}
