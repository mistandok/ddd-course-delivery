package postgre

import (
	"context"

	"delivery/internal/adapters/out/postgre/courier_repo"
	"delivery/internal/adapters/out/postgre/order_repo"
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
	db          *sqlx.DB
	trManager   *manager.Manager
	txGetter    TxGetter
	orderRepo   ports.OrderRepo
	courierRepo ports.CourierRepo
}

func NewUnitOfWork(
	db *sqlx.DB,
	trManager *manager.Manager,
	txGetter TxGetter,
) ports.UnitOfWork {
	uow := &UnitOfWork{}

	orderRepo := order_repo.NewRepository(db, txGetter)
	courierRepo := courier_repo.NewRepository(db, txGetter)

	uow.orderRepo = orderRepo
	uow.courierRepo = courierRepo
	uow.txGetter = txGetter
	uow.trManager = trManager
	uow.db = db

	return uow
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

func (u *UnitOfWork) CourierRepo() ports.CourierRepo {
	return u.courierRepo
}
