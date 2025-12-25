package kafka

import (
	"context"
	"fmt"

	"delivery/internal/adapters/in/kafka/common"
	"delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/generated/queues/basketpb"

	"github.com/google/uuid"
)

type BasketConfirmedEventHandler struct {
	createOrderHandler create_order.CreateOrderHandler
}

func (h *BasketConfirmedEventHandler) Handle(ctx context.Context, event basketpb.BasketConfirmedIntegrationEvent) error {
	orderID := uuid.New()

	cmd, err := create_order.NewCreateOrderCommand(
		orderID,
		event.Address.Street,
		int64(event.Volume),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", common.ErrIncorrectMessage, err)
	}

	return h.createOrderHandler.Handle(ctx, cmd)
}

func NewBasketConfirmedEventHandler(createOrderHandler create_order.CreateOrderHandler) *BasketConfirmedEventHandler {
	return &BasketConfirmedEventHandler{createOrderHandler: createOrderHandler}
}
