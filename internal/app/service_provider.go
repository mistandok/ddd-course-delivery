package app

import (
	"log"

	httpv1 "delivery/internal/adapters/in/http/v1"
	"delivery/internal/adapters/in/kafka"
	kafkaConsumerCommon "delivery/internal/adapters/in/kafka/common"
	"delivery/internal/adapters/out/grpc/geo"
	kafkaProducerCommon "delivery/internal/adapters/out/kafka/common"
	"delivery/internal/adapters/out/kafka/mapper"
	"delivery/internal/adapters/out/postgre"
	"delivery/internal/adapters/out/postgre/courier_repo"
	"delivery/internal/adapters/out/postgre/event_publisher"
	"delivery/internal/adapters/out/postgre/order_repo"
	"delivery/internal/config"
	"delivery/internal/config/env"
	eventHandlers "delivery/internal/core/application/event_handlers"
	"delivery/internal/core/application/usecases/commands/add_storage_place"
	"delivery/internal/core/application/usecases/commands/assign_order"
	"delivery/internal/core/application/usecases/commands/create_courier"
	"delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/core/application/usecases/commands/move_couriers_and_complete_order"
	"delivery/internal/core/application/usecases/queries/get_all_couriers"
	"delivery/internal/core/application/usecases/queries/get_all_uncompleted_orders"
	"delivery/internal/core/domain/model/event"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/crons"
	"delivery/internal/generated/queues/basketpb"
	"delivery/internal/generated/queues/orderpb"
	"delivery/internal/pkg/closer"
	eventPublisher "delivery/internal/pkg/event_publisher"
	"delivery/internal/pkg/outbox"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
)

type serviceProvider struct {
	pgConfig    *config.PgConfig
	httpConfig  *config.HttpConfig
	geoConfig   *config.GeoConfig
	kafkaConfig *config.KafkaConfig
	db          *sqlx.DB
	trManager   *manager.Manager
	uowFactory  ports.UnitOfWorkFactory
	orderRepo   ports.OrderRepo
	courierRepo ports.CourierRepo

	// External clients
	geoClient ports.GeoClient

	// HTTP
	httpHandlers *httpv1.DeliveryService

	// Cron Jobs
	moveCouriersJob cron.Job
	assignOrdersJob cron.Job

	// Kafka Consumers
	basketConfirmedConsumerGroup *kafkaConsumerCommon.KafkaConsumer[*basketpb.BasketConfirmedIntegrationEvent]
	basketConfirmedEventHandler  *kafka.BasketConfirmedEventHandler

	// Kafka Producers
	orderCreatedProducer                  ports.EventProducer[*event.OrderCreated]
	orderCompletedProducer                ports.EventProducer[*event.OrderCompleted]
	fromOrderCreatedToIntegrationMapper   *mapper.OrderCreatedMapper
	fromOrderCompletedToIntegrationMapper *mapper.OrderCompletedMapper

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

	// Event Handlers
	orderCreatedHandler   *eventHandlers.OrderCreatedHandler
	orderCompletedHandler *eventHandlers.OrderCompletedHandler

	// Event Publishers
	eventPublisher ports.EventPublisher

	// Event Registry
	eventRegistry outbox.EventRegistry
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
		s.orderRepo = order_repo.NewRepository(s.DB(), trmsqlx.DefaultCtxGetter, s.EventPublisher())
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
		s.uowFactory = postgre.NewUnitOfWorkFactory(s.DB(), s.TRManager(), trmsqlx.DefaultCtxGetter, s.EventPublisher())
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
		s.createOrderHandler = create_order.NewCreateOrderHandler(s.UOWFactory(), s.GeoClient())
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

func (s *serviceProvider) HttpConfig() *config.HttpConfig {
	if s.httpConfig == nil {
		httpConfig, err := config.NewHttpConfigSearcher().Get()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = httpConfig
	}

	return s.httpConfig
}

func (s *serviceProvider) GeoConfig() *config.GeoConfig {
	if s.geoConfig == nil {
		geoConfig, err := config.NewGeoConfigSearcher().Get()
		if err != nil {
			log.Fatalf("failed to get geo config: %v", err)
		}

		s.geoConfig = geoConfig
	}

	return s.geoConfig
}

func (s *serviceProvider) HttpHandlers() *httpv1.DeliveryService {
	if s.httpHandlers == nil {
		s.httpHandlers = httpv1.NewDeliveryService(
			s.GetAllCouriersHandler(),
			s.CreateCourierHandler(),
			s.GetAllUncompletedOrdersHandler(),
			s.CreateOrderHandler(),
		)
	}

	return s.httpHandlers
}

// Cron Jobs

func (s *serviceProvider) MoveCouriersJob() cron.Job {
	if s.moveCouriersJob == nil {
		job, err := crons.NewMoveCouriersJob(s.MoveCouriersAndCompleteOrderHandler())
		if err != nil {
			log.Fatalf("cannot create MoveCouriersJob: %v", err)
		}
		s.moveCouriersJob = job
	}

	return s.moveCouriersJob
}

func (s *serviceProvider) AssignOrdersJob() cron.Job {
	if s.assignOrdersJob == nil {
		job, err := crons.NewAssignOrdersJob(s.AssignOrderHandler())
		if err != nil {
			log.Fatalf("cannot create AssignOrdersJob: %v", err)
		}
		s.assignOrdersJob = job
	}

	return s.assignOrdersJob
}

// External Clients

func (s *serviceProvider) GeoClient() ports.GeoClient {
	if s.geoClient == nil {
		geoHost := s.GeoConfig().Address()
		client, closerFunc := geo.NewGeoClient(geoHost)

		closer.Add(closerFunc)
		s.geoClient = client
	}

	return s.geoClient
}

func (s *serviceProvider) KafkaConfig() *config.KafkaConfig {
	if s.kafkaConfig == nil {
		kafkaConfig, err := env.NewKafkaCfgSearcher().Get()
		if err != nil {
			log.Fatalf("failed to get kafka config: %v", err)
		}

		s.kafkaConfig = kafkaConfig
	}
	return s.kafkaConfig
}

func (s *serviceProvider) BasketConfirmedEventHandler() *kafka.BasketConfirmedEventHandler {
	if s.basketConfirmedEventHandler == nil {
		s.basketConfirmedEventHandler = kafka.NewBasketConfirmedEventHandler(s.CreateOrderHandler())
	}

	return s.basketConfirmedEventHandler
}

func (s *serviceProvider) BasketConfirmedConsumerGroup() *kafkaConsumerCommon.KafkaConsumer[*basketpb.BasketConfirmedIntegrationEvent] {
	if s.basketConfirmedConsumerGroup == nil {
		consumerGroup, err := kafkaConsumerCommon.NewKafkaConsumerGroup[*basketpb.BasketConfirmedIntegrationEvent](
			[]string{s.KafkaConfig().Host},
			s.KafkaConfig().ConsumerGroup,
			s.KafkaConfig().BasketConfirmedTopic,
			s.BasketConfirmedEventHandler(),
		)
		if err != nil {
			log.Fatalf("failed to create basket confirmed consumer group: %v", err)
		}

		s.basketConfirmedConsumerGroup = consumerGroup
	}

	return s.basketConfirmedConsumerGroup
}

func (s *serviceProvider) OrderCreatedHandler() *eventHandlers.OrderCreatedHandler {
	if s.orderCreatedHandler == nil {
		s.orderCreatedHandler = eventHandlers.NewOrderCreatedHandler(s.OrderCreatedProducer())
	}
	return s.orderCreatedHandler
}

func (s *serviceProvider) OrderCompletedHandler() *eventHandlers.OrderCompletedHandler {
	if s.orderCompletedHandler == nil {
		s.orderCompletedHandler = eventHandlers.NewOrderCompletedHandler(s.OrderCompletedProducer())
	}
	return s.orderCompletedHandler
}

func (s *serviceProvider) EventPublisher() eventPublisher.EventPublisher {
	if s.eventPublisher == nil {
		s.eventPublisher = event_publisher.NewEventPublisher(trmsqlx.DefaultCtxGetter, s.EventRegistry())
	}
	return s.eventPublisher
}

func (s *serviceProvider) FromOrderCreatedToIntegrationMapper() *mapper.OrderCreatedMapper {
	if s.fromOrderCreatedToIntegrationMapper == nil {
		s.fromOrderCreatedToIntegrationMapper = mapper.NewOrderCreatedMapper()
	}
	return s.fromOrderCreatedToIntegrationMapper
}

func (s *serviceProvider) FromOrderCompletedToIntegrationMapper() *mapper.OrderCompletedMapper {
	if s.fromOrderCompletedToIntegrationMapper == nil {
		s.fromOrderCompletedToIntegrationMapper = mapper.NewOrderCompletedMapper()
	}
	return s.fromOrderCompletedToIntegrationMapper
}

func (s *serviceProvider) OrderCreatedProducer() ports.EventProducer[*event.OrderCreated] {
	if s.orderCreatedProducer == nil {
		producer, err := kafkaProducerCommon.NewKafkaProducer[
			*event.OrderCreated,
			*orderpb.OrderCreatedIntegrationEvent,
			*mapper.OrderCreatedMapper,
		](
			[]string{s.KafkaConfig().Host},
			s.KafkaConfig().OrderChangedTopic,
			s.FromOrderCreatedToIntegrationMapper(),
		)
		if err != nil {
			log.Fatalf("failed to create order created producer: %v", err)
		}
		s.orderCreatedProducer = producer
	}
	return s.orderCreatedProducer
}

func (s *serviceProvider) OrderCompletedProducer() ports.EventProducer[*event.OrderCompleted] {
	if s.orderCompletedProducer == nil {
		producer, err := kafkaProducerCommon.NewKafkaProducer[
			*event.OrderCompleted,
			*orderpb.OrderCompletedIntegrationEvent,
			*mapper.OrderCompletedMapper,
		](
			[]string{s.KafkaConfig().Host},
			s.KafkaConfig().OrderChangedTopic,
			s.FromOrderCompletedToIntegrationMapper(),
		)
		if err != nil {
			log.Fatalf("failed to create order completed producer: %v", err)
		}
		s.orderCompletedProducer = producer
	}
	return s.orderCompletedProducer
}

func (s *serviceProvider) EventRegistry() outbox.EventRegistry {
	if s.eventRegistry == nil {
		eventRegistry, err := outbox.NewEventRegistry()
		if err != nil {
			log.Fatalf("failed to create event registry: %v", err)
		}
		s.eventRegistry = eventRegistry
	}

	return s.eventRegistry
}
