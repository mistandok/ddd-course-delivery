package order_repo

import (
	"context"
	modelOrder "delivery/internal/core/domain/model/order"
)

func (r *Repository) publishDomainEvents(ctx context.Context, order *modelOrder.Order) error {
	events := order.DomainEvents()

	for _, e := range events {
		err := r.eventPublisher.Publish(ctx, e)
		if err != nil {
			return err
		}
	}

	return nil
}
