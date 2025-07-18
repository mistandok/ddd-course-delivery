package shared_kernel

import (
	"delivery/internal/pkg/errs"
	"fmt"
	"math/rand/v2"
)

const (
	minX, maxX = 1, 10
	minY, maxY = 1, 10
)

type location struct {
	x     int64
	y     int64
	isSet bool
}

func NewLocation(x int64, y int64) (location, error) {
	if x < minX || x > maxX {
		return location{}, errs.NewValueIsOutOfRangeError("x", x, minX, maxX)
	}

	if y < minY || y > maxY {
		return location{}, errs.NewValueIsOutOfRangeError("y", y, minY, maxY)
	}

	return location{x: x, y: y, isSet: true}, nil
}

func NewRandomLocation() (location, error) {
	x, err := randomInt64InRange(minX, maxX)
	if err != nil {
		return location{}, err
	}

	y, err := randomInt64InRange(minY, maxY)
	if err != nil {
		return location{}, err
	}

	return NewLocation(x, y)
}

func (l location) X() int64 {
	return l.x
}

func (l location) Y() int64 {
	return l.y
}

func (l location) IsSet() bool {
	return l.isSet
}

func (l location) Equals(other location) bool {
	return l == other
}

func (l location) DistanceTo(other location) int64 {
	return abs(l.x-other.x) + abs(l.y-other.y)
}

// TODO: move to pkg if it will be used in other places
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// TODO: move to pkg if it will be used in other places
func randomInt64InRange(min, max int64) (int64, error) {
	if min > max {
		return 0, errs.NewValueIsInvalidErrorWithCause("min", fmt.Errorf("min is greater than max"))
	}
	return min + rand.Int64N(max-min+1), nil
}
