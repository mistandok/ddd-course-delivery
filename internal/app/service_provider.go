package app

import (
	"delivery/internal/adapters/out/postgre"
	"delivery/internal/config"
	"delivery/internal/config/env"
	"delivery/internal/core/ports"
	"delivery/pkg/closer"
	"log"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

type serviceProvider struct {
	pgConfig  *config.PgConfig
	db        *sqlx.DB
	trManager *manager.Manager
	uow       ports.UnitOfWork
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

func (s *serviceProvider) UOW() ports.UnitOfWork {
	if s.uow == nil {
		s.uow = postgre.NewUnitOfWork(s.DB(), s.TRManager())
	}

	return s.uow
}
