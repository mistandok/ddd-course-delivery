package courier_repo

import (
	"context"
	"database/sql"
	"errors"

	modelCourier "delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, courier *modelCourier.Courier) error {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	courierDTO, storagePlacesDTO := DomainToDTO(courier)

	courierExists, err := r.courierExists(ctx, tx, courierDTO.ID)
	if err != nil {
		return err
	}

	if !courierExists {
		return errs.NewObjectNotFoundError("courier", courierDTO.ID)
	}

	err = r.updateCourierWithOptimisticLock(ctx, tx, courierDTO)
	if err != nil {
		return err
	}

	err = r.updateStoragePlaces(ctx, tx, courierDTO.ID, storagePlacesDTO)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) updateCourierWithOptimisticLock(ctx context.Context, tx trmsqlx.Tr, courierDTO *CourierDTO) error {
	query, args, err := squirrel.Update("courier").
		Where(squirrel.Eq{"id": courierDTO.ID}).
		Where(squirrel.Eq{"version": courierDTO.Version}).
		Set("name", courierDTO.Name).
		Set("speed", courierDTO.Speed).
		Set("location", squirrel.Expr("POINT(?, ?)", courierDTO.Location.X, courierDTO.Location.Y)).
		Set("version", courierDTO.Version+1).
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
			return errs.NewVersionIsInvalidError("courier", errors.New("version mismatch"))
		}
		return err
	}

	return nil
}

func (r *Repository) updateStoragePlaces(ctx context.Context, tx trmsqlx.Tr, courierID uuid.UUID, storagePlacesDTO []StoragePlaceDTO) error {
	err := r.deleteStoragePlaces(ctx, tx, courierID)
	if err != nil {
		return err
	}

	err = r.insertStoragePlaces(ctx, tx, storagePlacesDTO)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) deleteStoragePlaces(ctx context.Context, tx trmsqlx.Tr, courierID uuid.UUID) error {
	deleteQuery, deleteArgs, err := squirrel.Delete("storage_place").
		Where(squirrel.Eq{"courier_id": courierID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, deleteQuery, deleteArgs...)
	return err
}

func (r *Repository) insertStoragePlaces(ctx context.Context, tx trmsqlx.Tr, storagePlacesDTO []StoragePlaceDTO) error {
	for _, spDTO := range storagePlacesDTO {
		spQuery, spArgs, err := squirrel.Insert("storage_place").
			Columns("id", "order_id", "courier_id", "volume", "name").
			Values(spDTO.ID, spDTO.OrderID, spDTO.CourierID, spDTO.Volume, spDTO.Name).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, spQuery, spArgs...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) courierExists(ctx context.Context, tx trmsqlx.Tr, id uuid.UUID) (bool, error) {
	checkQuery, checkArgs, err := squirrel.Select("1").
		From("courier").
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
