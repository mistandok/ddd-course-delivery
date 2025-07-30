package courier

import (
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	allowedVolume  int64  = 10
	zeroVolume     int64  = 0
	negativeVolume int64  = -1
	backpackName   string = "backpack"
)

func Test_Imposible_Create_Storage_Place_With_Invalid_Volume(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		volume int64
	}{
		{
			name:   "impossible create storage place with zero volume",
			volume: zeroVolume,
		},
		{
			name:   "impossible create storage place with negative volume",
			volume: negativeVolume,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewStoragePlace("test", test.volume)
			assert.ErrorIs(t, err, errs.ErrValueIsInvalid)
		})
	}
}

func Test_Imposible_Create_Storage_Place_With_Empty_Name(t *testing.T) {
	t.Parallel()
	// Arrange
	emptyName := ""

	// Act
	_, err := NewStoragePlace(emptyName, allowedVolume)

	// Assert
	assert.ErrorIs(t, err, errs.ErrValueIsInvalid)
}

func Test_Create_Storage_Place(t *testing.T) {
	t.Parallel()

	// Act
	storagePlace, err := NewStoragePlace(backpackName, allowedVolume)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, storagePlace.ID())
	assert.Nil(t, storagePlace.OrderID())
	assert.Equal(t, backpackName, storagePlace.Name())
	assert.Equal(t, allowedVolume, storagePlace.TotalVolume())
}

func Test_Imposible_Store_Order_In_Storage_Place_When_Order_Volume_Is_Greater_Than_Storage_Place_Volume(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)
	orderID := uuid.New()
	orderVolume := allowedVolume + 1

	// Act
	err := storagePlace.Store(orderID, orderVolume)

	// Assert
	assert.ErrorIs(t, err, errs.ErrValueIsInvalid)
}

func Test_Storage_Place_Is_Occupied_If_Successfull_Store_Order(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)
	orderID := uuid.New()
	orderVolume := allowedVolume

	// Act
	_ = storagePlace.Store(orderID, orderVolume)
	isOccupied := storagePlace.IsOccupied()

	// Assert
	assert.True(t, isOccupied)
}

func Test_Storage_Place_Is_Not_Occupied_If_Not_Store_Order(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)

	// Act
	isOccupied := storagePlace.IsOccupied()

	// Assert
	assert.False(t, isOccupied)
}

func Test_Impossible_Store_Order_In_Storage_Place_When_Storage_Place_Is_Occupied(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)
	firstOrderID := uuid.New()
	firstOrderVolume := allowedVolume
	secondOrderID := uuid.New()
	secondOrderVolume := allowedVolume

	// Act
	_ = storagePlace.Store(firstOrderID, firstOrderVolume)
	err := storagePlace.Store(secondOrderID, secondOrderVolume)

	// Assert
	assert.ErrorIs(t, err, errs.ErrValueIsInvalid)
}
