package postgre

import (
	"context"
	"delivery/internal/core/ports"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

var _ ports.UnitOfWork = (*UnitOfWork)(nil)

type TxGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type UnitOfWork struct {
	db        *sqlx.DB
	trManager *manager.Manager
	txGetter  TxGetter
	orderRepo ports.OrderRepo
}

func NewUnitOfWork(db *sqlx.DB, trManager *manager.Manager, txGetter TxGetter, orderRepo ports.OrderRepo) ports.UnitOfWork {
	return &UnitOfWork{
		db:        db,
		trManager: trManager,
		txGetter:  txGetter,
		orderRepo: orderRepo,
	}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return u.trManager.Do(ctx, fn)
}

func (u *UnitOfWork) DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr {
	return u.txGetter.DefaultTrOrDB(ctx, db)
}

func (u *UnitOfWork) OrderRepo() ports.OrderRepo {
	return u.orderRepo
}
