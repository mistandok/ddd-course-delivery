package postgre

import (
	"context"
	"delivery/internal/adapters/out/postgre/order_repo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
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

func Test_OrderRepoShouldGetOrder(t *testing.T) {
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

func Test_OrderRepoShouldUpdateOrder(t *testing.T) {
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Add(ctx, order)
	})

	// Act
	// Ничего не обновляем в заказе, просто пытаемся его обновить как есть
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Update(ctx, order)
	})

	// Assert
	assert.NoError(t, err)
}

func Test_ImpossibleToUpdateOrderWhenItNotExists(t *testing.T) {
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)

	// Act
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Update(ctx, order)
	})

	// Assert
	assert.ErrorIs(t, err, errs.ErrObjectNotFound)
}

func Test_ImpossobleToUpdateOrderWhenSomeoneElseUpdatedIt(t *testing.T) {
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	// Добавляем заказ
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Add(ctx, order)
	})
	// Обновляем заказ (предположим, что это сделал другой поток)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Update(ctx, order)
	})

	// Act
	// Пытаемся обновить заказ, который обновили в другом треде
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.OrderRepo().Update(ctx, order)
	})

	// Assert
	assert.ErrorIs(t, err, errs.ErrVersionIsInvalid)
}

func Test_OrderRepoShouldGetFirstInCreatedStatus(t *testing.T) {
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	oldestOrder, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	youngestOrder, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		_ = uow.OrderRepo().Add(ctx, oldestOrder)
		_ = uow.OrderRepo().Add(ctx, youngestOrder)

		return nil
	})

	// Act
	gettedOrder, err := uow.OrderRepo().GetFirstInCreatedStatus(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, oldestOrder.ID(), gettedOrder.ID())
}
