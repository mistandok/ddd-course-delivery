package order_repo

import (
	modelOrder "delivery/internal/core/domain/model/order"
)

func DomainToDTO(order *modelOrder.Order) *OrderDTO {
	return &OrderDTO{
		ID:        order.ID(),
		CourierID: order.CourierID(),
		Location: LocationDTO{
			X: order.Location().X(),
			Y: order.Location().Y(),
		},
		Volume:  order.Volume(),
		Status:  order.Status().String(),
		Version: order.Version(),
	}
}
