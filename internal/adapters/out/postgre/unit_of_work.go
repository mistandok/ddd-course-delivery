package postgre

import (
	"context"
	"delivery/internal/core/ports"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

var _ ports.UnitOfWork = (*UnitOfWork)(nil)

type UnitOfWork struct {
	db        *sqlx.DB
	trManager *manager.Manager
}

func NewUnitOfWork(db *sqlx.DB, trManager *manager.Manager) ports.UnitOfWork {
	return &UnitOfWork{db: db, trManager: trManager}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return u.trManager.Do(ctx, fn)
}
