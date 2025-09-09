package get_all_couriers

import (
	"context"

	"delivery/internal/pkg/errs"

	"database/sql"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type GetAllCouriersHandler interface {
	Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

var _ GetAllCouriersHandler = (*getAllCouriersHandler)(nil)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type getAllCouriersHandler struct {
	db       *sqlx.DB
	txGetter txGetter
}

func NewGetAllCouriersHandler(db *sqlx.DB, txGetter txGetter) *getAllCouriersHandler {
	return &getAllCouriersHandler{db: db, txGetter: txGetter}
}

func (h *getAllCouriersHandler) Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, errs.NewQueryIsInvalidError(query.QueryName())
	}

	tx := h.txGetter.DefaultTrOrDB(ctx, h.db)

	qry, args, err := squirrel.Select("id", "name", "location").
		From("courier").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return GetAllCouriersResponse{}, err
	}

	var couriers []CourierDTO
	err = tx.SelectContext(ctx, &couriers, qry, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetAllCouriersResponse{}, nil
		}
		return GetAllCouriersResponse{}, err
	}

	return GetAllCouriersResponse{
		Couriers: couriers,
	}, nil
}
