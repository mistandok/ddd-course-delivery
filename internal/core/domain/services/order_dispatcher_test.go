package services

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	aggCourier "delivery/internal/core/domain/model/courier"
	aggOrder "delivery/internal/core/domain/model/order"
	kernel "delivery/internal/core/domain/model/shared_kernel"
)

func TestCourierDispatcher_ImpossibleToDispatchOrderWithoutCouriers(t *testing.T) {
	// Arrange
	dispatcher := NewCourierDispatcher()
	order := getRandomOrder(t)
	couriers := []*aggCourier.Courier{}

	// Act
	_, err := dispatcher.Dispatch(order, couriers)

	// Assert
	assert.Error(t, err)
}

func TestCourierDispatcher_ImpossibleToDispatchMissingOrder(t *testing.T) {
	// Arrange
	dispatcher := NewCourierDispatcher()
	couriers := getRandomCouriers(t)

	// Act
	_, err := dispatcher.Dispatch(nil, couriers)

	// Assert
	assert.Error(t, err)
}

func TestCourierDispatcher_ImpossibleToDispatchOrderNotInCreationState(t *testing.T) {
	// Arrange
	dispatcher := NewCourierDispatcher()
	order := getRandomAssignedOrder(t)
	couriers := getRandomCouriers(t)
	// Act
	_, err := dispatcher.Dispatch(order, couriers)

	// Assert
	assert.Error(t, err)
}

func TestCourierDispatcher_SelectTheBestCourier(t *testing.T) {
	// Arrange
	dispatcher := NewCourierDispatcher()

	locationForCourierAndOrder, _ := kernel.NewLocation(1, 1)
	locationForCourierWhichIsFarFromOrder, _ := kernel.NewLocation(5, 5)

	order := getOrderWithLocation(t, locationForCourierAndOrder)

	expectedCourier := getCourierWithLocation(t, "courier-1", locationForCourierAndOrder)
	courierWhichIsFarFromOrder := getCourierWithLocation(t, "courier-2", locationForCourierWhichIsFarFromOrder)
	couriers := []*aggCourier.Courier{courierWhichIsFarFromOrder, expectedCourier}

	// Act
	assignedCourier, err := dispatcher.Dispatch(order, couriers)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCourier.ID(), assignedCourier.ID())
	assert.True(t, order.Status().Equals(aggOrder.StatusAssigned))
	assert.NotNil(t, order.CourierID())
	assert.Equal(t, *order.CourierID(), assignedCourier.ID())
}

func getRandomOrder(t *testing.T) *aggOrder.Order {
	t.Helper()

	location, err := kernel.NewRandomLocation()
	if err != nil {
		t.Fatalf("failed to create random location: %v", err)
	}

	order, err := aggOrder.NewOrder(uuid.New(), location, 1)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	return order
}

func getRandomAssignedOrder(t *testing.T) *aggOrder.Order {
	t.Helper()

	order := getRandomOrder(t)
	courier := getRandomCourier(t, "courier-1")
	err := order.Assign(courier.ID())
	if err != nil {
		t.Fatalf("failed to assign order to courier: %v", err)
	}

	return order
}

func getRandomCouriers(t *testing.T) []*aggCourier.Courier {
	t.Helper()

	couriers := []*aggCourier.Courier{}

	for i := range 2 {
		couriers = append(couriers, getRandomCourier(t, fmt.Sprintf("courier-%d", i)))
	}

	return couriers
}

func getRandomCourier(t *testing.T, name string) *aggCourier.Courier {
	t.Helper()

	location, err := kernel.NewRandomLocation()
	if err != nil {
		t.Fatalf("failed to create random location: %v", err)
	}

	courier, err := aggCourier.NewCourier(name, 100, location)
	if err != nil {
		t.Fatalf("failed to create courier: %v", err)
	}

	return courier
}

func getCourierWithLocation(t *testing.T, name string, location kernel.Location) *aggCourier.Courier {
	t.Helper()

	courier, err := aggCourier.NewCourier(name, 100, location)
	if err != nil {
		t.Fatalf("failed to create courier: %v", err)
	}

	return courier
}

func getOrderWithLocation(t *testing.T, location kernel.Location) *aggOrder.Order {
	t.Helper()

	order, err := aggOrder.NewOrder(uuid.New(), location, 1)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	return order
}
