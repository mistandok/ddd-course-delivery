package courier_repo

import (
	"context"

	modelCourier "delivery/internal/core/domain/model/courier"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
)

func (r *Repository) GetAllCouriersWithOrders(ctx context.Context) ([]*modelCourier.Courier, error) {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	couriersDTO, err := r.getCouriersWithOrdersDTO(ctx, tx)
	if err != nil {
		return nil, err
	}

	if len(couriersDTO) == 0 {
		return []*modelCourier.Courier{}, nil
	}

	courierIDs := r.extractCourierIDs(couriersDTO)

	storagePlacesByCourier, err := r.getStoragePlacesByCourierIDs(ctx, tx, courierIDs)
	if err != nil {
		return nil, err
	}

	return r.convertToCouriersWithOrdersDomain(couriersDTO, storagePlacesByCourier)
}

func (r *Repository) getCouriersWithOrdersDTO(ctx context.Context, tx trmsqlx.Tr) ([]CourierDTO, error) {
	query, args, err := squirrel.Select("c.id", "c.name", "c.speed", "c.location", "c.version").
		From("courier c").
		Join("storage_place sp ON c.id = sp.courier_id").
		Where("sp.order_id IS NOT NULL").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var couriersDTO []CourierDTO
	err = tx.SelectContext(ctx, &couriersDTO, query, args...)
	if err != nil {
		return nil, err
	}

	return couriersDTO, nil
}

func (r *Repository) convertToCouriersWithOrdersDomain(couriersDTO []CourierDTO, storagePlacesByCourier map[uuid.UUID][]StoragePlaceDTO) ([]*modelCourier.Courier, error) {
	var result []*modelCourier.Courier
	for _, courierDTO := range couriersDTO {
		storagePlaces := storagePlacesByCourier[courierDTO.ID]

		courier, err := DTOToDomain(&courierDTO, storagePlaces)
		if err != nil {
			return nil, err
		}

		result = append(result, courier)
	}

	return result, nil
}
