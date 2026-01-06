package outbox

import (
	"encoding/json"
	"testing"

	"delivery/internal/core/domain/model/event"
	"delivery/internal/pkg/outbox/mappers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EventRegistry_EncodeDomainEvent_SuccessfullyEncodesOrderCreated(t *testing.T) {
	// Arrange
	registry := setupRegistryWithOrderCreated(t)
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	domainEvent := event.NewOrderCreated(orderID)

	// Act
	message, err := registry.EncodeDomainEvent(domainEvent)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, domainEvent.GetID(), message.ID)
	assert.Equal(t, "OrderCreated", message.Name)
	assert.NotEmpty(t, message.Payload)
	assert.NotZero(t, message.OccurredAtUtc)
	assert.Nil(t, message.ProcessedAtUtc)

	var payloadData map[string]interface{}
	err = json.Unmarshal(message.Payload, &payloadData)
	require.NoError(t, err)
	assert.Equal(t, domainEvent.GetID().String(), payloadData["id"])
	assert.Equal(t, "order_created", payloadData["name"])
	assert.Equal(t, orderID.String(), payloadData["order_id"])
}

func Test_EventRegistry_EncodeDomainEvent_SuccessfullyEncodesOrderCompleted(t *testing.T) {
	// Arrange
	registry := setupRegistryWithOrderCompleted(t)
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	domainEvent := event.NewOrderCompleted(orderID)

	// Act
	message, err := registry.EncodeDomainEvent(domainEvent)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, domainEvent.GetID(), message.ID)
	assert.Equal(t, "OrderCompleted", message.Name)
	assert.NotEmpty(t, message.Payload)
	assert.NotZero(t, message.OccurredAtUtc)
	assert.Nil(t, message.ProcessedAtUtc)

	var payloadData map[string]interface{}
	err = json.Unmarshal(message.Payload, &payloadData)
	require.NoError(t, err)
	assert.Equal(t, domainEvent.GetID().String(), payloadData["id"])
	assert.Equal(t, "order_completed", payloadData["name"])
	assert.Equal(t, orderID.String(), payloadData["order_id"])
}

func Test_EventRegistry_EncodeDomainEvent_ReturnsErrorForUnregisteredEventType(t *testing.T) {
	// Arrange
	registry := setupEmptyRegistry(t)
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	domainEvent := event.NewOrderCreated(orderID)

	// Act
	message, err := registry.EncodeDomainEvent(domainEvent)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown domain event type: OrderCreated")
	assert.Equal(t, Message{}, message)
}

func Test_EventRegistry_EncodeDomainEvent_RoundTripPreservesEventFields(t *testing.T) {
	// Arrange
	registry := setupRegistryWithOrderCreated(t)
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	originalEvent := event.NewOrderCreated(orderID)

	// Act - Encode
	message, err := registry.EncodeDomainEvent(originalEvent)
	require.NoError(t, err)

	// Act - Decode
	decodedEvent, err := registry.DecodeDomainEvent(&message)
	require.NoError(t, err)

	// Assert
	decodedOrderCreated, ok := decodedEvent.(*event.OrderCreated)
	require.True(t, ok, "decoded event should be *OrderCreated")
	assert.Equal(t, originalEvent.GetID(), decodedOrderCreated.GetID())
	assert.Equal(t, originalEvent.GetName(), decodedOrderCreated.GetName())
	assert.Equal(t, originalEvent.GetOrderID(), decodedOrderCreated.GetOrderID())
}

// Helper functions

func setupEmptyRegistry(t *testing.T) EventRegistry {
	t.Helper()
	registry, err := NewEventRegistry()
	require.NoError(t, err)
	return registry
}

func setupRegistryWithOrderCreated(t *testing.T) EventRegistry {
	t.Helper()
	registry := setupEmptyRegistry(t)
	toMapper := mappers.NewOrderCreatedToJSONMapper()
	fromMapper := mappers.NewOrderCreatedFromJSONMapper()
	err := registry.RegisterDomainEvent(&event.OrderCreated{}, toMapper, fromMapper)
	require.NoError(t, err)
	return registry
}

func setupRegistryWithOrderCompleted(t *testing.T) EventRegistry {
	t.Helper()
	registry := setupEmptyRegistry(t)
	toMapper := mappers.NewOrderCompletedToJSONMapper()
	fromMapper := mappers.NewOrderCompletedFromJSONMapper()
	err := registry.RegisterDomainEvent(&event.OrderCompleted{}, toMapper, fromMapper)
	require.NoError(t, err)
	return registry
}
