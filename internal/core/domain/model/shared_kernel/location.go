package shared_kernel

import (
	"delivery/internal/pkg/errs"
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
	x := randomInt64InRange(minX, maxX)
	y := randomInt64InRange(minY, maxY)
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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func randomInt64InRange(min, max int64) int64 {
	// функция локальная, если мы некорректно зададим диапазон, то, как мне кажется, лучше сразу падать с паникой, так как это прям критическая ошибка для приложения.
	if min > max {
		panic("min is greater than max")
	}
	return min + rand.Int64N(max-min+1)
}
