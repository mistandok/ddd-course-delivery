package courier

import (
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/pointer"
	"errors"

	"github.com/google/uuid"
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int64
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int64) (StoragePlace, error) {
	if name == "" {
		return StoragePlace{}, errs.NewValueIsInvalidError("name")
	}

	if totalVolume <= 0 {
		return StoragePlace{}, errs.NewValueIsInvalidError("totalVolume")
	}

	storagePlace := StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}

	return storagePlace, nil
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
