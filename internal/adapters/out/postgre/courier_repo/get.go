package courier_repo

import (
	"context"
	"database/sql"
	"errors"

	modelCourier "delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*modelCourier.Courier, error) {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	courierQuery, courierArgs, err := squirrel.Select("id", "name", "speed", "location", "version").
		From("courier").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	courierDTO := &CourierDTO{}

	err = tx.GetContext(ctx, courierDTO, courierQuery, courierArgs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewObjectNotFoundError("courier", id)
		}
		return nil, err
	}

	storagePlacesQuery, storagePlacesArgs, err := squirrel.Select("id", "order_id", "courier_id", "volume", "name").
		From("storage_place").
		Where(squirrel.Eq{"courier_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var storagePlacesDTO []StoragePlaceDTO
	err = tx.SelectContext(ctx, &storagePlacesDTO, storagePlacesQuery, storagePlacesArgs...)
	if err != nil {
		return nil, err
	}

	return DTOToDomain(courierDTO, storagePlacesDTO)
}
