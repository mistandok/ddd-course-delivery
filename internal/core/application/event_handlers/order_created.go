package event_handlers

import (
	"context"
	"delivery/internal/core/domain/model/event"
	"log"
)

type OrderCreatedHandler struct {
}

func NewOrderCreatedHandler() *OrderCreatedHandler {
	return &OrderCreatedHandler{}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, event *event.OrderCreated) error {
	log.Printf("Order created: %v", event)
	return nil
}
