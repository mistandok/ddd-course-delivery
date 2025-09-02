package postgre

import (
	"context"
	"delivery/internal/adapters/out/postgre/order_repo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/testcnts"
	"log"
	"os"
	"testing"

	modelOrder "delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/model/shared_kernel"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var dbURL string
var uow ports.UnitOfWork

func TestMain(m *testing.M) {
	ctx := context.Background()

	testcnts.SetupTestEnvironment()

	postgresContainer, containerDBURL, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)

	db, trManager := setupDbEntities(containerDBURL)
	defer db.Close()

	orderRepo := order_repo.NewRepository(db, trmsqlx.DefaultCtxGetter)
	uow = NewUnitOfWork(db, trManager, trmsqlx.DefaultCtxGetter, orderRepo)

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

func Test_OrderRepoShouldAddOrder(t *testing.T) {
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)

	// Act
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Add(ctx, order)
	})

	// Assert
	assert.NoError(t, err)
}

func TestOrderRepoShouldGetOrder(t *testing.T) {
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Add(ctx, order)
	})

	// Act
	gettedOrder, err := uow.OrderRepo().Get(context.Background(), order.ID())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, order.ID(), gettedOrder.ID())
	assert.Equal(t, order.Location(), gettedOrder.Location())
	assert.Equal(t, order.Volume(), gettedOrder.Volume())
	assert.Equal(t, order.Status(), gettedOrder.Status())
	assert.Equal(t, order.CourierID(), gettedOrder.CourierID())
	assert.Equal(t, order.Version(), gettedOrder.Version())
}
