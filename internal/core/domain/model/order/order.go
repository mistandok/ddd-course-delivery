package order

import (
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/pkg/errs"
	"errors"

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

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if err := o.switchToStatus(StatusAssigned); err != nil {
		return err
	}

	o.courierID = &courierID

	return nil
}

func (o *Order) Complete() error {
	if err := o.switchToStatus(StatusCompleted); err != nil {
		return err
	}

	return nil
}

func (o *Order) switchToStatus(status Status) error {
	statusTransition := map[Status]Status{
		StatusCreated:   StatusAssigned,
		StatusAssigned:  StatusCompleted,
		StatusCompleted: StatusCompleted,
	}

	allowedNextStatus, ok := statusTransition[o.status]
	if !ok {
		return errs.NewValueIsInvalidErrorWithCause("status", errors.New("из текущего статуса заказа нельзя перейти в статус "+status.String()))
	}

	if allowedNextStatus != status {
		return errs.NewValueIsInvalidErrorWithCause("status", errors.New("из текущего статуса заказа нельзя перейти в статус "+status.String()))
	}

	o.status = status

	return nil
}
