package ports

import (
	aggCourier "delivery/internal/core/domain/model/courier"
	aggOrder "delivery/internal/core/domain/model/order"
)

//go:generate mockery --name OrderDispatcher --with-expecter --exported
type OrderDispatcher interface {
	Dispatch(order *aggOrder.Order, couriers []*aggCourier.Courier) (*aggCourier.Courier, error)
}
