package order_repo

import (
	"context"
	"database/sql"
	"errors"

	modelOrder "delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*modelOrder.Order, error) {
	query, args, err := squirrel.Select("id", "courier_id", "location", "volume", "status", "version").
		From(`"order"`).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	orderDTO := &OrderDTO{}

	err = r.db.GetContext(ctx, orderDTO, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("order", id)
		}

		return nil, err
	}

	return DTOToDomain(orderDTO)
}
