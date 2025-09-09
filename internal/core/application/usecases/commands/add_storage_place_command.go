package commands

import (
	"errors"

	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type AddStoragePlaceCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int64

	isValid bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int64) (AddStoragePlaceCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidErrorWithCause("courierID", errors.New("courierID is required"))
	}

	if name == "" {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidErrorWithCause("name", errors.New("name is required"))
	}

	if totalVolume <= 0 {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidErrorWithCause("totalVolume", errors.New("totalVolume must be greater than 0"))
	}

	return AddStoragePlaceCommand{courierID: courierID, name: name, totalVolume: totalVolume, isValid: true}, nil
}

func (c AddStoragePlaceCommand) CommandName() string {
	return "AddStoragePlaceCommand"
}

func (c AddStoragePlaceCommand) IsValid() bool {
	return c.isValid
}

func (c AddStoragePlaceCommand) CourierID() uuid.UUID {
	return c.courierID
}

func (c AddStoragePlaceCommand) Name() string {
	return c.name
}

func (c AddStoragePlaceCommand) TotalVolume() int64 {
	return c.totalVolume
}
