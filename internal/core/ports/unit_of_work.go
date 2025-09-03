package ports

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
	OrderRepo() OrderRepo
	CourierRepo() CourierRepo
}
