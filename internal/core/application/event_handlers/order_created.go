package event_handlers

import (
	"context"
	"delivery/internal/core/domain/model/event"
	"delivery/internal/core/ports"
	"log"
)

type OrderCreatedHandler struct {
	producer ports.EventProducer[*event.OrderCreated]
}

func NewOrderCreatedHandler(producer ports.EventProducer[*event.OrderCreated]) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		producer: producer,
	}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, event *event.OrderCreated) error {
	log.Printf("Order created: %v", event)
	return h.producer.Publish(ctx, event)
}
