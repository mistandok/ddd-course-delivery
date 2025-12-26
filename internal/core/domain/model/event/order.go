package event

import (
	"delivery/internal/pkg/ddd"

	"github.com/google/uuid"
)

type EventName string

const (
	EventNameOrderCreated   EventName = "order_created"
	EventNameOrderCompleted EventName = "order_completed"
)

var _ ddd.DomainEvent = (*OrderCreated)(nil)
var _ ddd.DomainEvent = (*OrderCompleted)(nil)

type OrderCreated struct {
	id   uuid.UUID
	name EventName

	orderID uuid.UUID
}

func NewOrderCreated(orderID uuid.UUID) *OrderCreated {
	return &OrderCreated{
		id:      uuid.New(),
		name:    EventNameOrderCreated,
		orderID: orderID,
	}
}

func (e *OrderCreated) GetID() uuid.UUID {
	return e.id
}

func (e *OrderCreated) GetName() string {
	return string(e.name)
}

func (e *OrderCreated) GetOrderID() uuid.UUID {
	return e.orderID
}

type OrderCompleted struct {
	id   uuid.UUID
	name EventName

	orderID uuid.UUID
}

func NewOrderCompleted(orderID uuid.UUID) *OrderCompleted {
	return &OrderCompleted{
		id:      uuid.New(),
		name:    EventNameOrderCompleted,
		orderID: orderID,
	}
}

func (e *OrderCompleted) GetID() uuid.UUID {
	return e.id
}

func (e *OrderCompleted) GetName() string {
	return string(e.name)
}

func (e *OrderCompleted) GetOrderID() uuid.UUID {
	return e.orderID
}
