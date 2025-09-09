package get_all_couriers

import (
	"context"
	"delivery/internal/pkg/errs"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

type GetAllCouriersHandler interface {
	Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

var _ GetAllCouriersHandler = (*getAllCouriersHandler)(nil)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type getAllCouriersHandler struct {
	txGetter txGetter
}

func NewGetAllCouriersHandler(txGetter txGetter) *getAllCouriersHandler {
	return &getAllCouriersHandler{txGetter: txGetter}
}

func (h *getAllCouriersHandler) Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, errs.NewQueryIsInvalidError(query.QueryName())
	}

	tx := h.txGetter.DefaultTrOrDB(ctx, nil)

	sql, args, err := squirrel.Select("id", "name", "location").
		From("courier").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return GetAllCouriersResponse{}, err
	}

	var couriers []CourierDTO
	err = tx.SelectContext(ctx, &couriers, sql, args...)
	if err != nil {
		return GetAllCouriersResponse{}, err
	}

	return GetAllCouriersResponse{
		Couriers: couriers,
	}, nil
}
