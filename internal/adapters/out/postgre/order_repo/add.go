package order_repo

import (
	"context"

	modelOrder "delivery/internal/core/domain/model/order"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) Add(ctx context.Context, order *modelOrder.Order) error {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	orderDTO := DomainToDTO(order)

	query, args, err := squirrel.Insert(`"order"`).
		Columns("id", "courier_id", "location", "volume", "status", "version").
		Values(
			orderDTO.ID,
			orderDTO.CourierID,
			squirrel.Expr("POINT(?, ?)", orderDTO.Location.X, orderDTO.Location.Y),
			orderDTO.Volume,
			orderDTO.Status,
			orderDTO.Version,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return r.publishDomainEvents(ctx, tx, order)
}
