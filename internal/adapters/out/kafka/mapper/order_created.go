package mapper

import (
	"delivery/internal/adapters/out/kafka/common"
	"delivery/internal/core/domain/model/event"
	"delivery/internal/generated/queues/orderpb"
)

type OrderCreatedMapper struct {
}

func NewOrderCreatedMapper() *OrderCreatedMapper {
	return &OrderCreatedMapper{}
}

func (m *OrderCreatedMapper) Map(domainEvent *event.OrderCreated) common.IntegrationEvent[*orderpb.OrderCreatedIntegrationEvent] {
	event := &orderpb.OrderCreatedIntegrationEvent{
		OrderId: domainEvent.GetOrderID().String(),
	}

	return *common.NewIntegrationEvent[*orderpb.OrderCreatedIntegrationEvent](event, domainEvent.GetID().String())
}
