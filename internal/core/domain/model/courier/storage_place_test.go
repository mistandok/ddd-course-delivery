package courier

import (
	"testing"

	"delivery/internal/pkg/errs"

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

func Test_Storage_Place_Is_Occupied_If_Successful_Store_Order(t *testing.T) {
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

func Test_If_Storage_Place_Is_Not_Occupied_Then_Can_Clear_Without_Problems(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)
	orderID := uuid.New()

	// Act
	err := storagePlace.Clear(orderID)

	// Assert
	assert.NoError(t, err)
}

func Test_Impossible_Clear_Order_In_Storage_Place_When_Order_Is_Not_Stored(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)
	firstOrderID := uuid.New()
	secondOrderID := uuid.New()

	// Act
	_ = storagePlace.Store(firstOrderID, allowedVolume)
	err := storagePlace.Clear(secondOrderID)

	// Assert
	assert.ErrorIs(t, err, errs.ErrObjectNotFound)
}

func Test_Clearing_Order_In_Storage_Place_When_Order_Is_Stored(t *testing.T) {
	t.Parallel()

	// Arrange
	storagePlace, _ := NewStoragePlace(backpackName, allowedVolume)
	orderID := uuid.New()
	orderVolume := allowedVolume

	// Act
	_ = storagePlace.Store(orderID, orderVolume)
	err := storagePlace.Clear(orderID)

	// Assert
	assert.NoError(t, err)
}
