package courier_repo

import (
	modelCourier "delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/shared_kernel"
)

func DomainToDTO(courier *modelCourier.Courier) (*CourierDTO, []StoragePlaceDTO) {
	courierDTO := &CourierDTO{
		ID:   courier.ID(),
		Name: courier.Name(),
		Speed: courier.Speed(),
		Location: LocationDTO{
			X: courier.Location().X(),
			Y: courier.Location().Y(),
		},
		Version: courier.Version(),
	}

	storagePlaces := make([]StoragePlaceDTO, 0, len(courier.StoragePlaces()))
	
	for _, sp := range courier.StoragePlaces() {
		storagePlaces = append(storagePlaces, StoragePlaceDTO{
			ID:        sp.ID(),
			Name:      sp.Name(),
			Volume:    sp.TotalVolume(),
			OrderID:   sp.OrderID(),
			CourierID: courier.ID(),
		})
	}

	return courierDTO, storagePlaces
}

func DTOToDomain(courierDTO *CourierDTO, storagePlacesDTO []StoragePlaceDTO) (*modelCourier.Courier, error) {
	location, err := shared_kernel.NewLocation(courierDTO.Location.X, courierDTO.Location.Y)
	if err != nil {
		return nil, err
	}

	storagePlaces := make([]*modelCourier.StoragePlace, 0, len(storagePlacesDTO))
	
	for _, spDTO := range storagePlacesDTO {
		sp := modelCourier.LoadStoragePlaceFromRepo(spDTO.ID, spDTO.Name, spDTO.Volume, spDTO.OrderID)
		storagePlaces = append(storagePlaces, sp)
	}

	return modelCourier.LoadCourierFromRepo(
		courierDTO.ID,
		courierDTO.Name,
		courierDTO.Speed,
		location,
		storagePlaces,
		courierDTO.Version,
	), nil
}