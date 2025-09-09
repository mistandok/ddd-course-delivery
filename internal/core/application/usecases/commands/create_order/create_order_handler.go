package create_order

import (
	"context"
	"errors"

	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderHandler interface {
	Handle(ctx context.Context, command CreateOrderCommand) error
}

var _ CreateOrderHandler = (*createOrderHandler)(nil)

type createOrderHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateOrderHandler(uowFactory ports.UnitOfWorkFactory) CreateOrderHandler {
	return &createOrderHandler{uowFactory: uowFactory}
}

func (h *createOrderHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewCommandIsInvalidErrorWithCause(command.CommandName(), errors.New("should use NewCreateOrderCommand to create a command"))
	}

	uow := h.uowFactory.NewUOW()

	err := uow.Do(ctx, func(ctx context.Context) error {
		// TODO: потом перейдем на другой способ генерации локации
		randomLocation, uowErr := shared_kernel.NewRandomLocation()
		if uowErr != nil {
			return uowErr
		}

		order, uowErr := order.NewOrder(command.OrderID(), randomLocation, command.Volume())
		if uowErr != nil {
			return uowErr
		}

		uowErr = uow.OrderRepo().Add(ctx, order)
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
