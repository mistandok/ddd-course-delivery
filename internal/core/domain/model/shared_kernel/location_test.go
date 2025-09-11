package shared_kernel

import (
	"testing"

	"delivery/internal/pkg/errs"

	"github.com/stretchr/testify/assert"
)

func Test_Impossible_To_Create_Location_With_Coordinates_Out_Of_Range(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		x    int64
		y    int64
	}{
		{
			name: "x is out of min range",
			x:    0,
			y:    1,
		},
		{
			name: "y is out of min range",
			x:    1,
			y:    0,
		},
		{
			name: "x and y are out of min range",
			x:    0,
			y:    0,
		},
		{
			name: "x is out of max range",
			x:    11,
			y:    1,
		},
		{
			name: "y is out of max range",
			x:    1,
			y:    11,
		},
		{
			name: "x and y are out of max range",
			x:    11,
			y:    11,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			_, err := NewLocation(test.x, test.y)

			// Assert
			assert.ErrorIs(t, err, errs.ErrValueIsOutOfRange)
		})
	}
}

func Test_Locations_With_Same_Coordinates_Should_Be_Equal(t *testing.T) {
	// Arrange
	location, _ := NewLocation(1, 1)
	otherLocation, _ := NewLocation(1, 1)

	// Act
	isEqual := location.Equals(otherLocation)

	// Assert
	assert.True(t, isEqual)
}

func Test_Locations_With_Different_Coordinates_Should_Be_Different(t *testing.T) {
	tests := []struct {
		name   string
		x      int64
		y      int64
		otherX int64
		otherY int64
	}{
		{
			name:   "x is different",
			x:      1,
			y:      1,
			otherX: 2,
			otherY: 1,
		},
		{
			name:   "y is different",
			x:      1,
			y:      1,
			otherX: 1,
			otherY: 2,
		},
		{
			name:   "x and y are different",
			x:      1,
			y:      1,
			otherX: 2,
			otherY: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			location, _ := NewLocation(test.x, test.y)
			otherLocation, _ := NewLocation(test.otherX, test.otherY)

			// Act
			isEqual := location.Equals(otherLocation)

			// Assert
			assert.False(t, isEqual)
		})
	}
}

func Test_Correct_Distance_To_Other_Location(t *testing.T) {
	tests := []struct {
		name                                string
		courierAndFinalDestinationLocations func() (Location, Location)
		expectedDistance                    int64
	}{
		{
			name:                                "same location",
			courierAndFinalDestinationLocations: courierAndFinalDestinationIsSameLocation,
			expectedDistance:                    0,
		},
		{
			name:                                "beetwen courier and final destination is one step by x",
			courierAndFinalDestinationLocations: beetwenCourierAndFinalDestinationIsOneStepByX,
			expectedDistance:                    1,
		},
		{
			name:                                "beetwen courier and final destination is one step by y",
			courierAndFinalDestinationLocations: beetwenCourierAndFinalDestinationIsOneStepByY,
			expectedDistance:                    1,
		},
		{
			name:                                "beetwen courier and final destination is five steps by x and y in summary",
			courierAndFinalDestinationLocations: beetwenCourierAndFinalDestinationIsFiveStepsByXAndYInSummary,
			expectedDistance:                    5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			courierLocation, finalDestinationLocation := test.courierAndFinalDestinationLocations()

			// Act
			distance := courierLocation.DistanceTo(finalDestinationLocation)

			// Assert
			assert.Equal(t, test.expectedDistance, distance)
		})
	}
}

func Test_Random_Location_Should_Be_In_Range(t *testing.T) {
	// Act
	location, err := NewRandomLocation()

	// Assert
	assert.NoError(t, err)
	assert.True(t, location.X() >= minX && location.X() <= maxX)
	assert.True(t, location.Y() >= minY && location.Y() <= maxY)
}

func courierAndFinalDestinationIsSameLocation() (Location, Location) {
	courierLocation, _ := NewLocation(1, 1)
	finalDestinationLocation, _ := NewLocation(1, 1)

	return courierLocation, finalDestinationLocation
}

func beetwenCourierAndFinalDestinationIsOneStepByX() (Location, Location) {
	courierLocation, _ := NewLocation(1, 1)
	finalDestinationLocation, _ := NewLocation(2, 1)

	return courierLocation, finalDestinationLocation
}

func beetwenCourierAndFinalDestinationIsOneStepByY() (Location, Location) {
	courierLocation, _ := NewLocation(1, 1)
	finalDestinationLocation, _ := NewLocation(1, 2)

	return courierLocation, finalDestinationLocation
}

func beetwenCourierAndFinalDestinationIsFiveStepsByXAndYInSummary() (Location, Location) {
	courierLocation, _ := NewLocation(4, 9)
	finalDestinationLocation, _ := NewLocation(2, 6)

	return courierLocation, finalDestinationLocation
}
