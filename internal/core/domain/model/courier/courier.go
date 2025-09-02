package courier

import (
	"delivery/internal/core/domain/model/order"
	kernel "delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

const (
	defaultStoragePlaceName         = "Сумка"
	defaultStoragePlaceVolume int64 = 10
)

type Courier struct {
	id            uuid.UUID
	name          string
	speed         int64
	location      kernel.Location
	storagePlaces []*StoragePlace
	version       int64
}

func NewCourier(name string, speed int64, location kernel.Location) (*Courier, error) {
	storagePlace, err := NewStoragePlace(defaultStoragePlaceName, defaultStoragePlaceVolume)
	if err != nil {
		return nil, err
	}

	if speed <= 0 {
		return nil, errs.NewValueIsInvalidErrorWithCause("speed", errors.New("speed must be greater than 0"))
	}

	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: []*StoragePlace{storagePlace},
	}, nil
}

func LoadCourierFromRepo(id uuid.UUID, name string, speed int64, location kernel.Location, storagePlaces []*StoragePlace, version int64) *Courier {
	return &Courier{
		id:            id,
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: storagePlaces,
		version:       version,
	}
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}

	return c.id == other.id
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int64 {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlaces
}

func (c *Courier) Version() int64 {
	return c.version
}

func (c *Courier) AddStoragePlace(name string, volume int64) error {
	storagePlace, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, storagePlace)
	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) bool {
	if order == nil {
		return false
	}

	for _, storagePlace := range c.storagePlaces {
		if storagePlace.CanStore(order.Volume()) {
			return true
		}
	}
	return false
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsInvalidErrorWithCause("order", errors.New("order is nil"))
	}

	if !c.CanTakeOrder(order) {
		return errs.NewValueIsInvalidErrorWithCause("order", errors.New("courier has no storage place with enough volume"))
	}

	for _, storagePlace := range c.storagePlaces {
		if storagePlace.CanStore(int64(order.Volume())) {
			if err := storagePlace.Store(order.ID(), order.Volume()); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	storagePlace, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return err
	}

	if err := storagePlace.Clear(order.ID()); err != nil {
		return err
	}

	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) float64 {
	distance := c.location.DistanceTo(target)
	return float64(distance) / float64(c.speed)
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsSet() {
		return errs.NewValueIsRequiredError("target")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int64(dx)
	newY := c.location.Y() + int64(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation

	return nil

}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	for _, storagePlace := range c.storagePlaces {
		if storagePlace.OrderID() != nil && *storagePlace.OrderID() == orderID {
			return storagePlace, nil
		}
	}

	return nil, errs.NewObjectNotFoundError("storage place", orderID)
}
