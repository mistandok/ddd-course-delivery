package order_repo

import (
	"context"
	modelOrder "delivery/internal/core/domain/model/order"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*modelOrder.Order, error) {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	query, args, err := squirrel.Select("id", "courier_id", "location", "volume", "status", "version").
		From(`"order"`).
		Where(squirrel.Eq{"status": modelOrder.StatusAssigned.String()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var ordersDTO []OrderDTO
	err = tx.SelectContext(ctx, &ordersDTO, query, args...)
	if err != nil {
		return nil, err
	}

	var result []*modelOrder.Order
	for _, orderDTO := range ordersDTO {
		order, err := DTOToDomain(&orderDTO)
		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}

	return result, nil
}
