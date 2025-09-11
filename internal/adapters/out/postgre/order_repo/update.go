package order_repo

import (
	"context"
	"database/sql"
	"errors"

	modelOrder "delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, order *modelOrder.Order) error {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	orderDTO := DomainToDTO(order)

	orderExists, err := r.orderExists(ctx, tx, orderDTO.ID)
	if err != nil {
		return err
	}

	if !orderExists {
		return errs.NewObjectNotFoundError("order", orderDTO.ID)
	}

	query, args, err := squirrel.Update(`"order"`).
		Where(squirrel.Eq{"id": orderDTO.ID}).
		Where(squirrel.Eq{"version": orderDTO.Version}).
		Set("courier_id", orderDTO.CourierID).
		Set("location", squirrel.Expr("POINT(?, ?)", orderDTO.Location.X, orderDTO.Location.Y)).
		Set("volume", orderDTO.Volume).
		Set("status", orderDTO.Status).
		Set("version", orderDTO.Version+1).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	var id uuid.UUID

	row := tx.QueryRowContext(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewVersionIsInvalidError("order", errors.New("version mismatch"))
		}

		return err
	}

	return nil
}

func (r *Repository) orderExists(ctx context.Context, tx trmsqlx.Tr, id uuid.UUID) (bool, error) {
	// Сначала проверяем существование записи
	checkQuery, checkArgs, err := squirrel.Select("1").
		From(`"order"`).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = tx.GetContext(ctx, &exists, checkQuery, checkArgs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
