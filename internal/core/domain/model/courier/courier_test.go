package courier

import (
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/model/shared_kernel"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const ()

func Test_Create_Courier_With_Default_StoragePlace(t *testing.T) {
	// Arrange
	location, _ := shared_kernel.NewRandomLocation()

	// Act
	courier, err := NewCourier("John Doe", 10, location)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", courier.Name())
	assert.Equal(t, location, courier.Location())
	assert.Equal(t, 1, len(courier.StoragePlaces()))
	assert.Equal(t, defaultStoragePlaceName, courier.StoragePlaces()[0].Name())
	assert.Equal(t, defaultStoragePlaceVolume, courier.StoragePlaces()[0].TotalVolume())
}

func Test_Courier_Can_Add_New_Storage_Place(t *testing.T) {
	// Arrange
	courier := newCourier(t)

	// Act
	_ = courier.AddStoragePlace("Ящик", 20)

	// Assert
	assert.Equal(t, 2, len(courier.StoragePlaces()))
	assert.Equal(t, "Ящик", courier.StoragePlaces()[1].Name())
	assert.Equal(t, int64(20), courier.StoragePlaces()[1].TotalVolume())
}

func Test_Courier_Can_Take_Order_If_At_Least_One_Storage_Place_Has_Enough_Volume(t *testing.T) {
	// Arrange
	testCases := []struct {
		name         string
		courier      func() *Courier
		orderForTake func() *order.Order
	}{
		{
			name: "Courier has empty storage place",
			courier: func() *Courier {
				courier := newCourier(t)
				return courier
			},
			orderForTake: func() *order.Order {
				return newOrderWithRandomLocationAndSettedVolume(t, 5)
			},
		},
		{
			name: "Courier has storage place in another place",
			courier: func() *Courier {
				courier := newCourier(t)
				order := newOrderWithRandomLocationAndSettedVolume(t, 5)
				_ = courier.TakeOrder(order)

				_ = courier.AddStoragePlace("Ящик", 20)

				return courier
			},
			orderForTake: func() *order.Order {
				return newOrderWithRandomLocationAndSettedVolume(t, 5)
			},
		},
	}

	// Act
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			courier := testCase.courier()
			order := testCase.orderForTake()

			canTakeOrder := courier.CanTakeOrder(order)
			assert.True(t, canTakeOrder)
		})
	}
}

func Test_Courier_Can_Not_Take_Order_If_Storage_Place_Is_Occupied(t *testing.T) {
	// Arrange
	courier := newCourier(t)
	placedOrder := newOrderWithRandomLocationAndSettedVolume(t, 5)
	newOrder := newOrderWithRandomLocationAndSettedVolume(t, 5)

	// Act
	_ = courier.TakeOrder(placedOrder)
	canTakeOrder := courier.CanTakeOrder(newOrder)

	// Assert
	assert.False(t, canTakeOrder)
}

func Test_Courier_Impossible_To_Take_Order_If_Storage_Place_Is_Occupied_By_Another_Order(t *testing.T) {
	// Arrange
	courier := newCourier(t)
	placedOrder := newOrderWithRandomLocationAndSettedVolume(t, 5)
	newOrder := newOrderWithRandomLocationAndSettedVolume(t, 5)

	// Act
	_ = courier.TakeOrder(placedOrder)
	err := courier.TakeOrder(newOrder)

	// Assert
	assert.Error(t, err)
}

func Test_Courier_Take_New_Order(t *testing.T) {
	// Arrange
	courier := newCourier(t)
	order := newOrderWithRandomLocationAndSettedVolume(t, 5)

	// Act
	err := courier.TakeOrder(order)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, order.ID(), *courier.StoragePlaces()[0].OrderID())
}

func Test_Courier_Can_Complete_Order(t *testing.T) {
	// Arrange
	courier := newCourier(t)
	order := newOrderWithRandomLocationAndSettedVolume(t, 5)

	// Act
	_ = courier.TakeOrder(order)
	err := courier.CompleteOrder(order)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, courier.StoragePlaces()[0].OrderID())
}

func Test_Calculate_Time_To_Location(t *testing.T) {
	// Arrange
	startLocation, _ := shared_kernel.NewLocation(1, 1)
	targetLocation, _ := shared_kernel.NewLocation(5, 5)
	courier, _ := NewCourier("John Doe", 2, startLocation)

	// Act
	time, _ := courier.CalculateTimeToLocation(targetLocation)

	// Assert
	assert.GreaterOrEqual(t, time, 4.0)
}

func newCourier(t *testing.T) *Courier {
	t.Helper()

	location, _ := shared_kernel.NewRandomLocation()
	courier, err := NewCourier("John Doe", 10, location)
	if err != nil {
		t.Fatal(err)
	}

	return courier
}

func newOrderWithRandomLocationAndSettedVolume(t *testing.T, volume int64) *order.Order {
	t.Helper()

	location, _ := shared_kernel.NewRandomLocation()
	order, err := order.NewOrder(uuid.New(), location, volume)
	if err != nil {
		t.Fatal(err)
	}

	return order
}
