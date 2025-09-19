package ports

import "delivery/internal/core/domain/model/shared_kernel"

type GeoClient interface {
	GetGeolocation(street string) (shared_kernel.Location, error)
}
