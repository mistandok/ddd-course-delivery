package order

import (
	"testing"

	"delivery/internal/core/domain/model/event"
	"delivery/internal/core/domain/model/shared_kernel"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Create_Order_With_Valid_Parameters(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	location, _ := shared_kernel.NewRandomLocation()
	volume := int64(10)

	// Act
	order, err := NewOrder(orderID, location, volume)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orderID, order.ID())
	assert.Equal(t, location, order.Location())
	assert.Equal(t, volume, order.Volume())
	assert.Equal(t, StatusCreated, order.Status())
	assert.Nil(t, order.CourierID())
}

func Test_NewOrder_Raises_OrderCreated_Event(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	location, _ := shared_kernel.NewRandomLocation()
	volume := int64(10)

	// Act
	order, err := NewOrder(orderID, location, volume)

	// Assert
	assert.NoError(t, err)
	events := order.DomainEvents()
	assert.Len(t, events, 1)

	orderCreatedEvent, ok := events[0].(*event.OrderCreated)
	assert.True(t, ok, "Expected event to be *event.OrderCreated")
	assert.Equal(t, orderID, orderCreatedEvent.GetOrderID())
}

func Test_Cannot_Create_Order_With_Empty_OrderID(t *testing.T) {
	// Arrange
	location, _ := shared_kernel.NewRandomLocation()
	volume := int64(10)

	// Act
	order, err := NewOrder(uuid.Nil, location, volume)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}

func Test_Cannot_Create_Order_With_Empty_Location(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	location := shared_kernel.Location{}
	volume := int64(10)

	// Act
	order, err := NewOrder(orderID, location, volume)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}

func Test_Cannot_Create_Order_With_Zero_Volume(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	location, _ := shared_kernel.NewRandomLocation()
	volume := int64(0)

	// Act
	order, err := NewOrder(orderID, location, volume)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}

func Test_Cannot_Create_Order_With_Negative_Volume(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	location, _ := shared_kernel.NewRandomLocation()
	volume := int64(-5)

	// Act
	order, err := NewOrder(orderID, location, volume)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}

func Test_Assign_Courier_To_Created_Order(t *testing.T) {
	// Arrange
	order := newValidOrder(t)
	courierID := uuid.New()

	// Act
	err := order.Assign(courierID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StatusAssigned, order.Status())
	assert.Equal(t, courierID, *order.CourierID())
}

func Test_Cannot_Assign_Courier_To_Already_Assigned_Order(t *testing.T) {
	// Arrange
	order := newValidOrder(t)
	firstCourierID := uuid.New()
	secondCourierID := uuid.New()

	// Act
	_ = order.Assign(firstCourierID)
	err := order.Assign(secondCourierID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, StatusAssigned, order.Status())
	assert.Equal(t, firstCourierID, *order.CourierID())
}

func Test_Cannot_Assign_Courier_To_Completed_Order(t *testing.T) {
	// Arrange
	order := newValidOrder(t)
	firstCourierID := uuid.New()
	secondCourierID := uuid.New()

	// Act
	_ = order.Assign(firstCourierID)
	_ = order.Complete()
	err := order.Assign(secondCourierID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, StatusCompleted, order.Status())
	assert.Equal(t, firstCourierID, *order.CourierID())
}

func Test_Complete_Assigned_Order(t *testing.T) {
	// Arrange
	order := newValidOrder(t)
	courierID := uuid.New()

	// Act
	_ = order.Assign(courierID)
	err := order.Complete()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, order.Status())
	assert.Equal(t, courierID, *order.CourierID())
}

func Test_Complete_Raises_OrderCompleted_Event(t *testing.T) {
	// Arrange
	order := newValidOrder(t)
	courierID := uuid.New()
	orderID := order.ID()

	// Act
	_ = order.Assign(courierID)
	err := order.Complete()

	// Assert
	assert.NoError(t, err)
	events := order.DomainEvents()
	assert.Len(t, events, 2)

	orderCreatedEvent, ok := events[0].(*event.OrderCreated)
	assert.True(t, ok, "Expected first event to be *event.OrderCreated")
	assert.Equal(t, orderID, orderCreatedEvent.GetOrderID())

	orderCompletedEvent, ok := events[1].(*event.OrderCompleted)
	assert.True(t, ok, "Expected second event to be *event.OrderCompleted")
	assert.Equal(t, orderID, orderCompletedEvent.GetOrderID())
}

func Test_Cannot_Complete_Created_Order_Without_Assignment(t *testing.T) {
	// Arrange
	order := newValidOrder(t)

	// Act
	err := order.Complete()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, StatusCreated, order.Status())
	assert.Nil(t, order.CourierID())
}

func Test_Failed_Complete_Does_Not_Raise_OrderCompleted_Event(t *testing.T) {
	// Arrange
	order := newValidOrder(t)

	// Act
	err := order.Complete()

	// Assert
	assert.Error(t, err)
	events := order.DomainEvents()
	assert.Len(t, events, 1)

	_, ok := events[0].(*event.OrderCreated)
	assert.True(t, ok, "Expected event to be *event.OrderCreated")
}

func Test_Cannot_Complete_Already_Completed_Order(t *testing.T) {
	// Arrange
	order := newValidOrder(t)
	courierID := uuid.New()

	// Act
	_ = order.Assign(courierID)
	_ = order.Complete()
	err := order.Complete()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, StatusCompleted, order.Status())
}

func newValidOrder(t *testing.T) *Order {
	t.Helper()

	orderID := uuid.New()
	location, _ := shared_kernel.NewRandomLocation()
	volume := int64(10)

	order, err := NewOrder(orderID, location, volume)
	if err != nil {
		t.Fatal(err)
	}

	return order
}
