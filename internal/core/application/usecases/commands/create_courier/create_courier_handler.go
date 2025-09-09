package create_courier

import (
	"context"
	"errors"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierHandler interface {
	Handle(ctx context.Context, command CreateCourierCommand) error
}

var _ CreateCourierHandler = (*createCourierHandler)(nil)

type createCourierHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateCourierHandler(uowFactory ports.UnitOfWorkFactory) CreateCourierHandler {
	return &createCourierHandler{uowFactory: uowFactory}
}

func (h *createCourierHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewCommandIsInvalidErrorWithCause(command.CommandName(), errors.New("should use NewCreateCourierCommand to create a command"))
	}

	uow := h.uowFactory.NewUOW()

	err := uow.Do(ctx, func(ctx context.Context) error {
		// TODO: потом перейдем на другой способ генерации локации
		randomLocation, uowErr := shared_kernel.NewRandomLocation()
		if uowErr != nil {
			return uowErr
		}

		courier, uowErr := courier.NewCourier(command.Name(), command.Speed(), randomLocation)
		if uowErr != nil {
			return uowErr
		}

		uowErr = uow.CourierRepo().Add(ctx, courier)
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
