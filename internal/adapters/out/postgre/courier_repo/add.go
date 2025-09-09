package courier_repo

import (
	"context"

	modelCourier "delivery/internal/core/domain/model/courier"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) Add(ctx context.Context, courier *modelCourier.Courier) error {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	courierDTO, storagePlacesDTO := DomainToDTO(courier)

	courierQuery, courierArgs, err := squirrel.Insert("courier").
		Columns("id", "name", "speed", "location", "version").
		Values(
			courierDTO.ID,
			courierDTO.Name,
			courierDTO.Speed,
			squirrel.Expr("POINT(?, ?)", courierDTO.Location.X, courierDTO.Location.Y),
			courierDTO.Version,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, courierQuery, courierArgs...)
	if err != nil {
		return err
	}

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
