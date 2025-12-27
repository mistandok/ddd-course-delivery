package event_handlers

import (
	"context"
	"delivery/internal/core/domain/model/event"
	"delivery/internal/core/ports"
	"log"
)

type OrderCompletedHandler struct {
	producer ports.EventProducer[*event.OrderCompleted]
}

func NewOrderCompletedHandler(producer ports.EventProducer[*event.OrderCompleted]) *OrderCompletedHandler {
	return &OrderCompletedHandler{
		producer: producer,
	}
}

func (h *OrderCompletedHandler) Handle(ctx context.Context, event *event.OrderCompleted) error {
	log.Printf("Order completed: %v", event)
	return h.producer.Publish(ctx, event)
}
