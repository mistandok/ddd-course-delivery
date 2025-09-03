package order_repo

import (
	modelOrder "delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/model/shared_kernel"
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

func DTOToDomain(orderDTO *OrderDTO) (*modelOrder.Order, error) {
	location, err := shared_kernel.NewLocation(orderDTO.Location.X, orderDTO.Location.Y)
	if err != nil {
		return nil, err
	}

	status := modelOrder.Status(orderDTO.Status)

	return modelOrder.LoadOrderFromRepo(
		orderDTO.ID,
		orderDTO.CourierID,
		location,
		orderDTO.Volume,
		status,
		orderDTO.Version,
	)
}
