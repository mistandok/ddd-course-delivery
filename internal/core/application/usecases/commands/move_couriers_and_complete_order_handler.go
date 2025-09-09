package commands

import (
	"context"
	modelCourier "delivery/internal/core/domain/model/courier"
	modelOrder "delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type MoveCouriersAndCompleteOrderHandler interface {
	Handle(ctx context.Context, command MoveCouriersAndFinishOrderCommand) error
}

var _ MoveCouriersAndCompleteOrderHandler = (*moveCouriersAndCompleteOrderHandler)(nil)

type moveCouriersAndCompleteOrderHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewMoveCouriersAndCompleteOrderHandler(uowFactory ports.UnitOfWorkFactory) MoveCouriersAndCompleteOrderHandler {
	return &moveCouriersAndCompleteOrderHandler{uowFactory: uowFactory}
}

func (h *moveCouriersAndCompleteOrderHandler) Handle(ctx context.Context, command MoveCouriersAndFinishOrderCommand) error {
	if !command.IsValid() {
		return errs.NewCommandIsInvalidErrorWithCause(
			command.CommandName(),
			errors.New("should use NewMoveCouriersAndFinishOrderCommand to create a command"),
		)
	}

	uow := h.uowFactory.NewUOW()

	err := uow.Do(ctx, func(ctx context.Context) error {
		assignedOrders, uowErr := uow.OrderRepo().GetAllInAssignedStatus(ctx)
		if uowErr != nil {
			return uowErr
		}

		for _, order := range assignedOrders {
			courier, uowErr := uow.CourierRepo().Get(ctx, *order.CourierID())
			if uowErr != nil {
				return uowErr
			}

			if uowErr := h.moveCourierAndCompleteOrder(courier, order); uowErr != nil {
				return uowErr
			}

			if uowErr := uow.CourierRepo().Update(ctx, courier); uowErr != nil {
				return uowErr
			}
			if uowErr := uow.OrderRepo().Update(ctx, order); uowErr != nil {
				return uowErr
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *moveCouriersAndCompleteOrderHandler) moveCourierAndCompleteOrder(courier *modelCourier.Courier, order *modelOrder.Order) error {
	if err := courier.Move(order.Location()); err != nil {
		return err
	}

	if courier.Location().Equals(order.Location()) {
		if err := order.Complete(); err != nil {
			return err
		}

		if err := courier.CompleteOrder(order); err != nil {
			return err
		}
	}

	return nil
}
