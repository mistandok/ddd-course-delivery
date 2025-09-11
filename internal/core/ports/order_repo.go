package ports

import (
	"context"

	modelOrder "delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
)

//go:generate mockery --name OrderRepo --with-expecter --exported
type OrderRepo interface {
	Add(ctx context.Context, order *modelOrder.Order) error
	Update(ctx context.Context, order *modelOrder.Order) error
	Get(ctx context.Context, id uuid.UUID) (*modelOrder.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*modelOrder.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*modelOrder.Order, error)
}
