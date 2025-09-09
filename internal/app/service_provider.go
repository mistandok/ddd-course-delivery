package app

import (
	"log"

	"delivery/internal/adapters/out/postgre"
	"delivery/internal/adapters/out/postgre/courier_repo"
	"delivery/internal/adapters/out/postgre/order_repo"
	"delivery/internal/config"
	"delivery/internal/config/env"
	"delivery/internal/core/application/usecases/commands/add_storage_place"
	"delivery/internal/core/application/usecases/commands/assign_order"
	"delivery/internal/core/application/usecases/commands/create_courier"
	"delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/core/application/usecases/commands/move_couriers_and_complete_order"
	"delivery/internal/core/application/usecases/queries/get_all_couriers"
	"delivery/internal/core/application/usecases/queries/get_all_uncompleted_orders"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/closer"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

type serviceProvider struct {
	pgConfig    *config.PgConfig
	db          *sqlx.DB
	trManager   *manager.Manager
	uowFactory  ports.UnitOfWorkFactory
	orderRepo   ports.OrderRepo
	courierRepo ports.CourierRepo

	// Domain Services
	orderDispatcher ports.OrderDispatcher

	// Command Handlers
	createOrderHandler                  create_order.CreateOrderHandler
	createeCourierHandler               create_courier.CreateCourierHandler
	addStoragePlaceHandler              add_storage_place.AddStoragePlaceHandler
	assignOrderHandler                  assign_order.AssignedOrderHandler
	moveCouriersAndCompleteOrderHandler move_couriers_and_complete_order.MoveCouriersAndCompleteOrderHandler

	// Query Handlers
	getAllCouriersHandler          get_all_couriers.GetAllCouriersHandler
	getAllUncompletedOrdersHandler get_all_uncompleted_orders.GetAllUncompletedOrdersHandler
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

func (s *serviceProvider) UOWFactory() ports.UnitOfWorkFactory {
	if s.uowFactory == nil {
		s.uowFactory = postgre.NewUnitOfWorkFactory(s.DB(), s.TRManager(), trmsqlx.DefaultCtxGetter)
	}

	return s.uowFactory
}

// Domain Services

func (s *serviceProvider) OrderDispatcher() ports.OrderDispatcher {
	if s.orderDispatcher == nil {
		s.orderDispatcher = services.NewCourierDispatcher()
	}

	return s.orderDispatcher
}

// Command Handlers

func (s *serviceProvider) CreateOrderHandler() create_order.CreateOrderHandler {
	if s.createOrderHandler == nil {
		s.createOrderHandler = create_order.NewCreateOrderHandler(s.UOWFactory())
	}

	return s.createOrderHandler
}

func (s *serviceProvider) CreateCourierHandler() create_courier.CreateCourierHandler {
	if s.createeCourierHandler == nil {
		s.createeCourierHandler = create_courier.NewCreateCourierHandler(s.UOWFactory())
	}

	return s.createeCourierHandler
}

func (s *serviceProvider) AddStoragePlaceHandler() add_storage_place.AddStoragePlaceHandler {
	if s.addStoragePlaceHandler == nil {
		s.addStoragePlaceHandler = add_storage_place.NewAddStoragePlaceHandler(s.UOWFactory())
	}

	return s.addStoragePlaceHandler
}

func (s *serviceProvider) AssignOrderHandler() assign_order.AssignedOrderHandler {
	if s.assignOrderHandler == nil {
		s.assignOrderHandler = assign_order.NewAssignedOrderHandler(s.UOWFactory(), s.OrderDispatcher())
	}

	return s.assignOrderHandler
}

func (s *serviceProvider) MoveCouriersAndCompleteOrderHandler() move_couriers_and_complete_order.MoveCouriersAndCompleteOrderHandler {
	if s.moveCouriersAndCompleteOrderHandler == nil {
		s.moveCouriersAndCompleteOrderHandler = move_couriers_and_complete_order.NewMoveCouriersAndCompleteOrderHandler(s.UOWFactory())
	}

	return s.moveCouriersAndCompleteOrderHandler
}

// Query Handlers

func (s *serviceProvider) GetAllCouriersHandler() get_all_couriers.GetAllCouriersHandler {
	if s.getAllCouriersHandler == nil {
		s.getAllCouriersHandler = get_all_couriers.NewGetAllCouriersHandler(s.DB(), trmsqlx.DefaultCtxGetter)
	}

	return s.getAllCouriersHandler
}

func (s *serviceProvider) GetAllUncompletedOrdersHandler() get_all_uncompleted_orders.GetAllUncompletedOrdersHandler {
	if s.getAllUncompletedOrdersHandler == nil {
		s.getAllUncompletedOrdersHandler = get_all_uncompleted_orders.NewGetAllUncompletedOrdersHandler(s.DB(), trmsqlx.DefaultCtxGetter)
	}

	return s.getAllUncompletedOrdersHandler
}
