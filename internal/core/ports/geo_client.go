package ports

import "delivery/internal/core/domain/model/shared_kernel"

//go:generate mockery --name=GeoClient --output=mocks --outpkg=mocks

type GeoClient interface {
	GetGeolocation(street string) (shared_kernel.Location, error)
}
