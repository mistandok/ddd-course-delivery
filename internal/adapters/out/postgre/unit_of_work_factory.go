package postgre

import (
	"delivery/internal/core/ports"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

var _ ports.UnitOfWorkFactory = (*UnitOfWorkFactory)(nil)

type UnitOfWorkFactory struct {
	db             *sqlx.DB
	trManager      *manager.Manager
	txGetter       TxGetter
	eventPublisher EventPublisher
}

func NewUnitOfWorkFactory(db *sqlx.DB, trManager *manager.Manager, txGetter TxGetter, eventPublisher EventPublisher) ports.UnitOfWorkFactory {
	return &UnitOfWorkFactory{db: db, trManager: trManager, txGetter: txGetter, eventPublisher: eventPublisher}
}

func (f *UnitOfWorkFactory) NewUOW() ports.UnitOfWork {
	return NewUnitOfWork(f.db, f.trManager, f.txGetter, f.eventPublisher)
}
