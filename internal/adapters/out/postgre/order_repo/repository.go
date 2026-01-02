package order_repo

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

var _ ports.OrderRepo = (*Repository)(nil)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type eventPublisher interface {
	Publish(ctx context.Context, event ddd.DomainEvent) error
}

type Repository struct {
	db             *sqlx.DB
	txGetter       txGetter
	eventPublisher eventPublisher
}

func NewRepository(db *sqlx.DB, txGetter txGetter, eventPublisher eventPublisher) *Repository {
	return &Repository{
		db:             db,
		txGetter:       txGetter,
		eventPublisher: eventPublisher,
	}
}
