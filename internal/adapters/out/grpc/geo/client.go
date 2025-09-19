package geo

import (
	"context"
	"log"
	"time"

	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.GeoClient = &geoClient{}

type closer func() error

type geoClient struct {
	client  geopb.GeoClient
	timeout time.Duration
}

type Option func(*geoClient)

// WithTimeout sets the timeout for gRPC calls
func WithTimeout(timeout time.Duration) Option {
	return func(c *geoClient) {
		c.timeout = timeout
	}
}

func NewGeoClient(host string, opts ...Option) (*geoClient, closer) {
	// Establish insecure connection
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to geo service: %v", err)
	}

	client := &geoClient{
		client:  geopb.NewGeoClient(conn),
		timeout: 30 * time.Second, // default timeout
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	closer := func() error {
		return conn.Close()
	}

	return client, closer
}

func (c *geoClient) GetGeolocation(street string) (shared_kernel.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &geopb.GetGeolocationRequest{
		Street: street,
	}

	resp, err := c.client.GetGeolocation(ctx, req)
	if err != nil {
		return shared_kernel.Location{}, err
	}

	if resp.Location == nil {
		return shared_kernel.Location{}, nil
	}

	// Convert from protobuf Location (int32) to shared_kernel.Location (int64)
	return shared_kernel.NewLocation(int64(resp.Location.X), int64(resp.Location.Y))
}
