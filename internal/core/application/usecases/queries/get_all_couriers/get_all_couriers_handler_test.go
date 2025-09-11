package get_all_couriers

import (
	"context"
	"log"
	"os"
	"testing"

	"delivery/internal/adapters/out/postgre"
	"delivery/internal/core/application/usecases/commands/create_courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/testcnts"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var dbURL string
var uowFactory ports.UnitOfWorkFactory
var handler GetAllCouriersHandler
var createCourierHandler create_courier.CreateCourierHandler

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

	uowFactory = postgre.NewUnitOfWorkFactory(db, trManager, trmsqlx.DefaultCtxGetter)
	handler = NewGetAllCouriersHandler(db, trmsqlx.DefaultCtxGetter)
	createCourierHandler = create_courier.NewCreateCourierHandler(uowFactory)

	dbURL = containerDBURL

	os.Exit(m.Run())
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

func createValidQuery() GetAllCouriersQuery {
	return NewGetAllCouriersQuery()
}

func addCourierViaHandler(t *testing.T, name string, speed int64) {
	t.Helper()
	command, err := create_courier.NewCreateCourierCommand(name, speed)
	assert.NoError(t, err)

	err = createCourierHandler.Handle(context.Background(), command)
	assert.NoError(t, err)
}

func Test_GetAllCouriersHandler_Handle_ValidQuery(t *testing.T) {
	cleanupDB(t)

	// Arrange
	addCourierViaHandler(t, "Courier 1", 10)
	addCourierViaHandler(t, "Courier 2", 15)

	query := createValidQuery()

	// Act
	response, err := handler.Handle(context.Background(), query)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, response.Couriers, 2)

	// Проверяем, что курьеры присутствуют в ответе
	courierNames := make([]string, 0, len(response.Couriers))
	for _, c := range response.Couriers {
		courierNames = append(courierNames, c.Name)
	}
	assert.Contains(t, courierNames, "Courier 1")
	assert.Contains(t, courierNames, "Courier 2")
}

func Test_GetAllCouriersHandler_Handle_EmptyDatabase(t *testing.T) {
	cleanupDB(t)

	// Arrange
	query := createValidQuery()

	// Act
	response, err := handler.Handle(context.Background(), query)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, response.Couriers)
}

func Test_GetAllCouriersHandler_Handle_InvalidQuery(t *testing.T) {
	cleanupDB(t)

	// Arrange
	query := GetAllCouriersQuery{isValid: false}

	// Act
	response, err := handler.Handle(context.Background(), query)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, response.Couriers)
}
