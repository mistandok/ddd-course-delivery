package assign_order

import (
	"context"
	"errors"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AssignedOrderHandler interface {
	Handle(ctx context.Context, command AssignedOrderCommand) error
}

var _ AssignedOrderHandler = (*assignedOrderHandler)(nil)

type assignedOrderHandler struct {
	uowFactory      ports.UnitOfWorkFactory
	orderDispatcher ports.OrderDispatcher
}

func NewAssignedOrderHandler(uowFactory ports.UnitOfWorkFactory, orderDispatcher ports.OrderDispatcher) AssignedOrderHandler {
	return &assignedOrderHandler{
		uowFactory:      uowFactory,
		orderDispatcher: orderDispatcher,
	}
}

func (h *assignedOrderHandler) Handle(ctx context.Context, command AssignedOrderCommand) error {
	if !command.IsValid() {
		return errs.NewCommandIsInvalidErrorWithCause(command.CommandName(), errors.New("command is invalid"))
	}

	uow := h.uowFactory.NewUOW()

	err := uow.Do(ctx, func(ctx context.Context) error {
		couriers, uowErr := uow.CourierRepo().GetAllFreeCouriers(ctx)
		if uowErr != nil {
			return uowErr
		}

		order, uowErr := uow.OrderRepo().GetFirstInCreatedStatus(ctx)
		if uowErr != nil {
			return uowErr
		}

		selectedCourier, uowErr := h.orderDispatcher.Dispatch(order, couriers)
		if uowErr != nil {
			return uowErr
		}

		if uowErr := uow.OrderRepo().Update(ctx, order); uowErr != nil {
			return uowErr
		}

		if uowErr := uow.CourierRepo().Update(ctx, selectedCourier); uowErr != nil {
			return uowErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
