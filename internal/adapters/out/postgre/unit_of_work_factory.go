package postgre

import (
	"delivery/internal/core/ports"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

var _ ports.UnitOfWorkFactory = (*UnitOfWorkFactory)(nil)

type UnitOfWorkFactory struct {
	db        *sqlx.DB
	trManager *manager.Manager
	txGetter  TxGetter
}

func NewUnitOfWorkFactory(db *sqlx.DB, trManager *manager.Manager, txGetter TxGetter) ports.UnitOfWorkFactory {
	return &UnitOfWorkFactory{db: db, trManager: trManager, txGetter: txGetter}
}

func (f *UnitOfWorkFactory) NewUOW() (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db, f.trManager, f.txGetter), nil
}
