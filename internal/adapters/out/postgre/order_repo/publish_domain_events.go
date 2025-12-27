package order_repo

import (
	"context"
	"delivery/internal/core/domain/model/event"
	modelOrder "delivery/internal/core/domain/model/order"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

func (r *Repository) publishDomainEvents(ctx context.Context, _ trmsqlx.Tr, order *modelOrder.Order) error {
	events := order.DomainEvents()

	for _, e := range events {
		switch e := e.(type) {
		case *event.OrderCreated:
			err := r.eventPublisher.Publish(ctx, e)
			if err != nil {
				return err
			}
		case *event.OrderCompleted:
			err := r.eventPublisher.Publish(ctx, e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
