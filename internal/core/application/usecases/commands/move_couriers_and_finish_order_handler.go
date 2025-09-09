package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type MoveCouriersAndFinishOrderHandler interface {
	Handle(ctx context.Context, command MoveCouriersAndFinishOrderCommand) error
}

var _ MoveCouriersAndFinishOrderHandler = (*moveCouriersAndFinishOrderHandler)(nil)

type moveCouriersAndFinishOrderHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewMoveCouriersAndFinishOrderHandler(uowFactory ports.UnitOfWorkFactory) MoveCouriersAndFinishOrderHandler {
	return &moveCouriersAndFinishOrderHandler{uowFactory: uowFactory}
}

func (h *moveCouriersAndFinishOrderHandler) Handle(ctx context.Context, command MoveCouriersAndFinishOrderCommand) error {
	if !command.IsValid() {
		return errs.NewCommandIsInvalidErrorWithCause(
			command.CommandName(),
			errors.New("should use NewMoveCouriersAndFinishOrderCommand to create a command"),
		)
	}

	uow := h.uowFactory.NewUOW()

	err := uow.Do(ctx, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
