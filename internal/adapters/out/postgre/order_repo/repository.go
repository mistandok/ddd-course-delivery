package order_repo

import (
	"context"
	"delivery/internal/core/ports"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

var _ ports.OrderRepo = (*Repository)(nil)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type Repository struct {
	db       *sqlx.DB
	txGetter txGetter
}

func NewRepository(db *sqlx.DB, txGetter txGetter) *Repository {
	return &Repository{
		db:       db,
		txGetter: txGetter,
	}
}
