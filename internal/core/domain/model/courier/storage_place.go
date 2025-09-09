package courier

import (
	"errors"

	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/pointer"

	"github.com/google/uuid"
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int64
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int64) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsInvalidError("name")
	}

	if totalVolume <= 0 {
		return nil, errs.NewValueIsInvalidError("totalVolume")
	}

	storagePlace := &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}

	return storagePlace, nil
}

func LoadStoragePlaceFromRepo(id uuid.UUID, name string, totalVolume int64, orderID *uuid.UUID) *StoragePlace {
	return &StoragePlace{
		id:          id,
		name:        name,
		totalVolume: totalVolume,
		orderID:     orderID,
	}
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int64 {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int64) error {
	if !s.CanStore(volume) {
		return errs.NewValueIsInvalidErrorWithCause("order", errors.New("order volume is greater than storage place volume"))
	}

	s.orderID = pointer.New(orderID)

	return nil
}

func (s *StoragePlace) IsOccupied() bool {
	return s.orderID != nil
}

func (s *StoragePlace) CanStore(volume int64) bool {
	if s.IsOccupied() {
		return false
	}

	return volume <= s.totalVolume
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if !s.IsOccupied() {
		return nil
	}

	if *s.orderID != orderID {
		return errs.NewObjectNotFoundError("order", orderID)
	}

	s.orderID = nil

	return nil
}

func (s *StoragePlace) Equal(other *StoragePlace) bool {
	if other == nil {
		return false
	}

	return s.id == other.id
}
