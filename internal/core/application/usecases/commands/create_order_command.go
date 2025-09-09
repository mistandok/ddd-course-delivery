package commands

import (
	"errors"

	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string
	volume  int64

	isValid bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int64) (CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errs.NewValueIsInvalidErrorWithCause("orderID", errors.New("orderID is required"))
	}

	if street == "" {
		return CreateOrderCommand{}, errs.NewValueIsInvalidErrorWithCause("street", errors.New("street is required"))
	}

	if volume <= 0 {
		return CreateOrderCommand{}, errs.NewValueIsInvalidErrorWithCause("volume", errors.New("volume must be greater than 0"))
	}

	return CreateOrderCommand{orderID: orderID, street: street, volume: volume, isValid: true}, nil
}

func (c CreateOrderCommand) CommandName() string {
	return "CreateOrderCommand"
}

func (c CreateOrderCommand) IsValid() bool {
	return c.isValid
}

func (c CreateOrderCommand) OrderID() uuid.UUID {
	return c.orderID
}

func (c CreateOrderCommand) Street() string {
	return c.street
}

func (c CreateOrderCommand) Volume() int64 {
	return c.volume
}
