package commands

import (
	"context"
	"errors"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AddStoragePlaceHandler interface {
	Handle(ctx context.Context, command AddStoragePlaceCommand) error
}

var _ AddStoragePlaceHandler = (*addStoragePlaceHandler)(nil)

type addStoragePlaceHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewAddStoragePlaceHandler(uowFactory ports.UnitOfWorkFactory) AddStoragePlaceHandler {
	return &addStoragePlaceHandler{uowFactory: uowFactory}
}

func (h *addStoragePlaceHandler) Handle(ctx context.Context, command AddStoragePlaceCommand) error {
	if !command.IsValid() {
		return errs.NewCommandIsInvalidErrorWithCause(command.CommandName(), errors.New("should use NewAddStoragePlaceCommand to create a command"))
	}

	uow := h.uowFactory.NewUOW()

	err := uow.Do(ctx, func(ctx context.Context) error {
		courier, uowErr := uow.CourierRepo().Get(ctx, command.CourierID())
		if uowErr != nil {
			return uowErr
		}

		err := courier.AddStoragePlace(command.Name(), command.TotalVolume())
		if err != nil {
			return err
		}

		uowErr = uow.CourierRepo().Update(ctx, courier)
		if uowErr != nil {
			return uowErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
