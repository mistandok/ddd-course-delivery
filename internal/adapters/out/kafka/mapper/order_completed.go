package mapper

import (
	"delivery/internal/adapters/out/kafka/common"
	"delivery/internal/core/domain/model/event"
	"delivery/internal/generated/queues/orderpb"
)

type OrderCompletedMapper struct {
}

func NewOrderCompletedMapper() *OrderCompletedMapper {
	return &OrderCompletedMapper{}
}

func (m *OrderCompletedMapper) Map(domainEvent *event.OrderCompleted) *common.IntegrationEvent[*orderpb.OrderCompletedIntegrationEvent] {
	event := &orderpb.OrderCompletedIntegrationEvent{
		OrderId: domainEvent.GetOrderID().String(),
	}

	return common.NewIntegrationEvent[*orderpb.OrderCompletedIntegrationEvent](event, domainEvent.GetID().String())
}
