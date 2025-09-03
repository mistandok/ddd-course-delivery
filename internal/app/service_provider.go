package app

import (
	"delivery/internal/adapters/out/postgre"
	"delivery/internal/adapters/out/postgre/courier_repo"
	"delivery/internal/adapters/out/postgre/order_repo"
	"delivery/internal/config"
	"delivery/internal/config/env"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/closer"
	"log"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

type serviceProvider struct {
	pgConfig    *config.PgConfig
	db          *sqlx.DB
	trManager   *manager.Manager
	uow         ports.UnitOfWork
	orderRepo   ports.OrderRepo
	courierRepo ports.CourierRepo
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PgConfig() *config.PgConfig {
	if s.pgConfig == nil {
		pgConfig, err := env.NewPgCfgSearcher().Get()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = pgConfig
	}

	return s.pgConfig
}

func (s *serviceProvider) DB() *sqlx.DB {
	if s.db == nil {
		db, err := sqlx.Connect("postgres", s.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}

		closer.Add(func() error {
			return db.Close()
		})

		s.db = db
	}

	return s.db
}

func (s *serviceProvider) TRManager() *manager.Manager {
	if s.trManager == nil {
		s.trManager = manager.Must(trmsqlx.NewDefaultFactory(s.DB()))
	}

	return s.trManager
}

func (s *serviceProvider) OrderRepo() ports.OrderRepo {
	if s.orderRepo == nil {
		s.orderRepo = order_repo.NewRepository(s.DB(), trmsqlx.DefaultCtxGetter)
	}

	return s.orderRepo
}

func (s *serviceProvider) CourierRepo() ports.CourierRepo {
	if s.courierRepo == nil {
		s.courierRepo = courier_repo.NewRepository(s.DB(), trmsqlx.DefaultCtxGetter)
	}

	return s.courierRepo
}

func (s *serviceProvider) UOW() ports.UnitOfWork {
	if s.uow == nil {
		s.uow = postgre.NewUnitOfWork(s.DB(), s.TRManager(), trmsqlx.DefaultCtxGetter)
	}

	return s.uow
}
