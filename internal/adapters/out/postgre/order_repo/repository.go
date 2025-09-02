package order_repo

import (
	"context"
	"delivery/internal/core/ports"

	"github.com/Masterminds/squirrel"

	modelOrder "delivery/internal/core/domain/model/order"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
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

func (r *Repository) Add(ctx context.Context, order *modelOrder.Order) error {
	tx := r.txGetter.DefaultTrOrDB(ctx, r.db)

	orderDTO := DomainToDTO(order)

	sql, args, err := squirrel.Insert("order").
		Columns("id", "courier_id", "location", "volume", "status", "version").
		Values(
			orderDTO.ID,
			orderDTO.CourierID,
			orderDTO.Location.String(),
			orderDTO.Volume,
			orderDTO.Status,
			orderDTO.Version,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, order *modelOrder.Order) error {
	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*modelOrder.Order, error) {
	return nil, nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*modelOrder.Order, error) {
	return nil, nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*modelOrder.Order, error) {
	return nil, nil
}
