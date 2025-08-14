package order

import (
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  shared_kernel.Location
	volume    int64
	status    Status
}

func NewOrder(orderID uuid.UUID, location shared_kernel.Location, volume int64) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if !location.IsSet() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		id:       orderID,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Location() shared_kernel.Location {
	return o.location
}

func (o *Order) Volume() int64 {
	return o.volume
}
