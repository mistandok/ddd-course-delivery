package ports

import (
	"context"

	modelCourier "delivery/internal/core/domain/model/courier"

	"github.com/google/uuid"
)

//go:generate mockery --name CourierRepo --with-expecter --exported
type CourierRepo interface {
	Add(ctx context.Context, courier *modelCourier.Courier) error
	Update(ctx context.Context, courier *modelCourier.Courier) error
	Get(ctx context.Context, id uuid.UUID) (*modelCourier.Courier, error)
	GetAllFreeCouriers(ctx context.Context) ([]*modelCourier.Courier, error)
	GetAllCouriersWithOrders(ctx context.Context) ([]*modelCourier.Courier, error)
}
