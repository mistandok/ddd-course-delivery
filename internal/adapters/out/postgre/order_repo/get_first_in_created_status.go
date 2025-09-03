package order_repo

import (
	"context"
	modelOrder "delivery/internal/core/domain/model/order"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*modelOrder.Order, error) {
	query, args, err := squirrel.Select("id", "courier_id", "location", "volume", "status", "version").
		From(`"order"`).
		Where(squirrel.Eq{"status": modelOrder.StatusCreated}).
		OrderBy("created_at").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	orderDTO := &OrderDTO{}

	err = r.db.GetContext(ctx, orderDTO, query, args...)
	if err != nil {
		return nil, err
	}

	return DTOToDomain(orderDTO)
}
