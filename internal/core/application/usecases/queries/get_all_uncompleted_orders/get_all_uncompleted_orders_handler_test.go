package get_all_uncompleted_orders

import (
	"context"
	"log"
	"os"
	"testing"

	"delivery/internal/adapters/out/postgre"
	"delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports"
	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/testcnts"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var dbURL string
var uowFactory ports.UnitOfWorkFactory
var handler GetAllUncompletedOrdersHandler
var createOrderHandler create_order.CreateOrderHandler
var geoClient ports.GeoClient

type fakeEventPublisher struct {
}

func (f *fakeEventPublisher) Publish(ctx context.Context, event ddd.DomainEvent) error {
	return nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	testcnts.SetupTestEnvironment()

	postgresContainer, containerDBURL, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate postgres container: %v", err)
		}
	}()

	db, trManager := setupDbEntities(containerDBURL)
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	eventPublisher := &fakeEventPublisher{}
	uowFactory = postgre.NewUnitOfWorkFactory(db, trManager, trmsqlx.DefaultCtxGetter, eventPublisher)
	handler = NewGetAllUncompletedOrdersHandler(db, trmsqlx.DefaultCtxGetter)

	// Setup mock GeoClient for integration tests
	mockGeoClient := setupMockGeoClient()
	geoClient = mockGeoClient
	createOrderHandler = create_order.NewCreateOrderHandler(uowFactory, geoClient)

	dbURL = containerDBURL

	os.Exit(m.Run())
}

func setupMockGeoClient() *mocks.GeoClient {
	mockGeoClient := &mocks.GeoClient{}
	// Return a fixed location for any street in integration tests
	location, _ := shared_kernel.NewLocation(5, 5)
	mockGeoClient.On("GetGeolocation", mock.Anything).Return(location, nil)
	return mockGeoClient
}

func setupDbEntities(dbURL string) (*sqlx.DB, *manager.Manager) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	return db, trManager
}

func cleanupDB(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		db, err := sqlx.Connect("postgres", dbURL)
		if err != nil {
			t.Fatalf("failed to connect to db for cleanup: %v", err)
		}
		defer db.Close()

		_, err = db.Exec("TRUNCATE TABLE storage_place, \"order\", courier RESTART IDENTITY CASCADE")
		if err != nil {
			t.Fatalf("failed to cleanup database: %v", err)
		}
	})
}

func createValidQuery() GetAllUncompletedOrdersQuery {
	return NewGetAllUncompletedOrdersQuery()
}

func addOrderViaHandler(t *testing.T, orderID uuid.UUID, street string, volume int64) {
	t.Helper()
	command, err := create_order.NewCreateOrderCommand(orderID, street, volume)
	assert.NoError(t, err)

	err = createOrderHandler.Handle(context.Background(), command)
	assert.NoError(t, err)
}

func Test_GetAllUncompletedOrdersHandler_Handle_ValidQuery(t *testing.T) {
	cleanupDB(t)

	// Arrange
	orderID1 := uuid.New()
	orderID2 := uuid.New()
	addOrderViaHandler(t, orderID1, "Street 1", 10)
	addOrderViaHandler(t, orderID2, "Street 2", 15)

	query := createValidQuery()

	// Act
	response, err := handler.Handle(context.Background(), query)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, response.Orders, 2)

	// Проверяем, что заказы присутствуют в ответе
	orderIDs := make([]uuid.UUID, 0, len(response.Orders))
	for _, o := range response.Orders {
		orderIDs = append(orderIDs, o.ID)
	}
	assert.Contains(t, orderIDs, orderID1)
	assert.Contains(t, orderIDs, orderID2)
}

func Test_GetAllUncompletedOrdersHandler_Handle_EmptyDatabase(t *testing.T) {
	cleanupDB(t)

	// Arrange
	query := createValidQuery()

	// Act
	response, err := handler.Handle(context.Background(), query)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, response.Orders)
}

func Test_GetAllUncompletedOrdersHandler_Handle_InvalidQuery(t *testing.T) {
	cleanupDB(t)

	// Arrange
	query := GetAllUncompletedOrdersQuery{isValid: false}

	// Act
	response, err := handler.Handle(context.Background(), query)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, response.Orders)
}
