package mappers

import (
	"encoding/json"
	"fmt"

	"delivery/internal/core/domain/model/event"
	"delivery/internal/pkg/ddd"

	"github.com/google/uuid"
)

var _ ToJSONMapper = &OrderCreatedToJSONMapper{}
var _ FromJSONMapper = &OrderCreatedFromJSONMapper{}

type OrderCreatedToJSONMapper struct{}

func NewOrderCreatedToJSONMapper() *OrderCreatedToJSONMapper {
	return &OrderCreatedToJSONMapper{}
}

func (m *OrderCreatedToJSONMapper) ToJSON(domainEvent ddd.DomainEvent) ([]byte, error) {
	e, ok := domainEvent.(*event.OrderCreated)
	if !ok {
		return nil, fmt.Errorf("expected *event.OrderCreated, got %T", domainEvent)
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

type OrderCreatedFromJSONMapper struct{}

func NewOrderCreatedFromJSONMapper() *OrderCreatedFromJSONMapper {
	return &OrderCreatedFromJSONMapper{}
}

func (m *OrderCreatedFromJSONMapper) FromJSON(data []byte) (ddd.DomainEvent, error) {
	var dto struct {
		ID      uuid.UUID        `json:"id"`
		Name    event.EventName  `json:"name"`
		OrderID uuid.UUID        `json:"order_id"`
	}

	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("unmarshal OrderCreated: %w", err)
	}

	return event.LoadOrderCreated(dto.ID, dto.Name, dto.OrderID), nil
}
