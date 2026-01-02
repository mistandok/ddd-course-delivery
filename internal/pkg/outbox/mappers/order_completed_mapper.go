package mappers

import (
	"encoding/json"
	"fmt"

	"delivery/internal/core/domain/model/event"
	"delivery/internal/pkg/ddd"

	"github.com/google/uuid"
)

var _ ToJSONMapper = &OrderCompletedToJSONMapper{}
var _ FromJSONMapper = &OrderCompletedFromJSONMapper{}

type OrderCompletedToJSONMapper struct{}

func NewOrderCompletedToJSONMapper() *OrderCompletedToJSONMapper {
	return &OrderCompletedToJSONMapper{}
}

func (m *OrderCompletedToJSONMapper) ToJSON(domainEvent ddd.DomainEvent) ([]byte, error) {
	e, ok := domainEvent.(*event.OrderCompleted)
	if !ok {
		return nil, fmt.Errorf("expected *event.OrderCompleted, got %T", domainEvent)
	}

	dto := struct {
		ID      uuid.UUID        `json:"id"`
		Name    event.EventName  `json:"name"`
		OrderID uuid.UUID        `json:"order_id"`
	}{
		ID:      e.GetID(),
		Name:    event.EventName(e.GetName()),
		OrderID: e.GetOrderID(),
	}

	return json.Marshal(dto)
}

type OrderCompletedFromJSONMapper struct{}

func NewOrderCompletedFromJSONMapper() *OrderCompletedFromJSONMapper {
	return &OrderCompletedFromJSONMapper{}
}

func (m *OrderCompletedFromJSONMapper) FromJSON(data []byte) (ddd.DomainEvent, error) {
	var dto struct {
		ID      uuid.UUID        `json:"id"`
		Name    event.EventName  `json:"name"`
		OrderID uuid.UUID        `json:"order_id"`
	}

	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("unmarshal OrderCompleted: %w", err)
	}

	return event.LoadOrderCompleted(dto.ID, dto.Name, dto.OrderID), nil
}
