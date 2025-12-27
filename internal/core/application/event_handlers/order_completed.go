package event_handlers

import (
	"context"
	"delivery/internal/core/domain/model/event"
	"log"
)

type OrderCompletedHandler struct {
}

func NewOrderCompletedHandler() *OrderCompletedHandler {
	return &OrderCompletedHandler{}
}

func (h *OrderCompletedHandler) Handle(ctx context.Context, event *event.OrderCompleted) error {
	log.Printf("Order completed: %v", event)
	return nil
}
