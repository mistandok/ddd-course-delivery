package get_all_uncompleted_orders

import (
	"context"

	"delivery/internal/pkg/errs"

	"database/sql"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type GetAllUncompletedOrdersHandler interface {
	Handle(ctx context.Context, query GetAllUncompletedOrdersQuery) (GetAllUncompletedOrdersResponse, error)
}

var _ GetAllUncompletedOrdersHandler = (*getAllUncompletedOrdersHandler)(nil)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type getAllUncompletedOrdersHandler struct {
	db       *sqlx.DB
	txGetter txGetter
}

func NewGetAllUncompletedOrdersHandler(db *sqlx.DB, txGetter txGetter) *getAllUncompletedOrdersHandler {
	return &getAllUncompletedOrdersHandler{db: db, txGetter: txGetter}
}

func (h *getAllUncompletedOrdersHandler) Handle(ctx context.Context, query GetAllUncompletedOrdersQuery) (GetAllUncompletedOrdersResponse, error) {
	if !query.IsValid() {
		return GetAllUncompletedOrdersResponse{}, errs.NewQueryIsInvalidError(query.QueryName())
	}

	tx := h.txGetter.DefaultTrOrDB(ctx, h.db)

	qry, args, err := squirrel.Select("id", "location").
		From("\"order\"").
		Where(squirrel.Or{
			squirrel.Eq{"status": "Assigned"},
			squirrel.Eq{"status": "Created"},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		if err == sql.ErrNoRows {
			return GetAllUncompletedOrdersResponse{}, nil
		}
		return GetAllUncompletedOrdersResponse{}, err
	}

	var orders []OrderDTO
	err = tx.SelectContext(ctx, &orders, qry, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetAllUncompletedOrdersResponse{}, nil
		}
		return GetAllUncompletedOrdersResponse{}, err
	}

	return GetAllUncompletedOrdersResponse{
		Orders: orders,
	}, nil
}
