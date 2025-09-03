package postgre

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/testcnts"
	"log"
	"os"
	"testing"

	modelCourier "delivery/internal/core/domain/model/courier"
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

	uow = NewUnitOfWork(db, trManager, trmsqlx.DefaultCtxGetter)

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

		// Очищаем таблицы в правильном порядке (из-за внешних ключей)
		_, err = db.Exec("TRUNCATE TABLE storage_place, \"order\", courier RESTART IDENTITY CASCADE")
		if err != nil {
			t.Fatalf("failed to cleanup database: %v", err)
		}
	})
}

func Test_OrderRepoShouldAddOrder(t *testing.T) {
	cleanupDB(t)
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
	cleanupDB(t)
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
	cleanupDB(t)
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
	cleanupDB(t)
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
	cleanupDB(t)
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
	cleanupDB(t)
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

func Test_OrderRepoShouldGetAllInAssignedStatus(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	assignedOrder, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	courier, _ := modelCourier.NewCourier("test", 10, randomLocation)
	_ = assignedOrder.Assign(courier.ID())
	// Добавляем курьера
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		_ = uow.CourierRepo().Add(ctx, courier)

		return nil
	})
	// Добавляем заказы
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		_ = uow.OrderRepo().Add(ctx, assignedOrder)
		_ = uow.OrderRepo().Add(ctx, order)

		return nil
	})

	// Act
	gettedOrders, err := uow.OrderRepo().GetAllInAssignedStatus(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gettedOrders))
	assert.Equal(t, assignedOrder.ID(), gettedOrders[0].ID())
}

func Test_CourierRepoShouldAddCourier(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	courier, _ := modelCourier.NewCourier("test", 10, randomLocation)

	// Act
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Add(ctx, courier)
	})

	// Assert
	assert.NoError(t, err)
}

func Test_CourierRepoShouldGetCourier(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	courier, _ := modelCourier.NewCourier("test", 10, randomLocation)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Add(ctx, courier)
	})

	// Act
	gettedCourier, err := uow.CourierRepo().Get(context.Background(), courier.ID())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, courier.ID(), gettedCourier.ID())
	assert.Equal(t, courier.Name(), gettedCourier.Name())
	assert.Equal(t, courier.Speed(), gettedCourier.Speed())
	assert.Equal(t, courier.Location(), gettedCourier.Location())
	assert.Equal(t, courier.Version(), gettedCourier.Version())
	assert.Equal(t, courier.StoragePlaces(), gettedCourier.StoragePlaces())
}

func Test_CourierRepoShouldUpdateCourier(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	courier, _ := modelCourier.NewCourier("test", 10, randomLocation)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Add(ctx, courier)
	})

	// Act
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Update(ctx, courier)
	})

	// Assert
	assert.NoError(t, err)
}

func Test_CourierRepoImpossibleToUpdateCourierWhenItNotExists(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	courier, _ := modelCourier.NewCourier("test", 10, randomLocation)

	// Act
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Update(ctx, courier)
	})

	// Assert
	assert.ErrorIs(t, err, errs.ErrObjectNotFound)
}

func Test_CourierRepoImpossibleToUpdateCourierWhenSomeoneElseUpdatedIt(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	courier, _ := modelCourier.NewCourier("test", 10, randomLocation)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Add(ctx, courier)
	})
	// Обновляем заказ (предположим, что это сделал другой поток)
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Update(ctx, courier)
	})

	// Act
	// Пытаемся обновить заказ, который обновили в другом треде
	err := uow.Do(context.Background(), func(ctx context.Context) error {
		return uow.CourierRepo().Update(ctx, courier)
	})

	// Assert
	assert.ErrorIs(t, err, errs.ErrVersionIsInvalid)
}

func Test_CourierRepoShouldGetAllFreeCouriers(t *testing.T) {
	cleanupDB(t)
	// Arrange
	randomLocation, _ := shared_kernel.NewRandomLocation()
	courierThatTakeOrder, _ := modelCourier.NewCourier("test", 10, randomLocation)
	freeCourier, _ := modelCourier.NewCourier("test", 10, randomLocation)
	order, _ := modelOrder.NewOrder(uuid.New(), randomLocation, 5)
	_ = courierThatTakeOrder.TakeOrder(order)

	// Добавляем заказ, свободного курьера и курьера, который взял заказ
	_ = uow.Do(context.Background(), func(ctx context.Context) error {
		_ = uow.OrderRepo().Add(ctx, order)
		_ = uow.CourierRepo().Add(ctx, freeCourier)
		_ = uow.CourierRepo().Add(ctx, courierThatTakeOrder)

		return nil
	})

	// Act
	gettedCouriers, err := uow.CourierRepo().GetAllFreeCouriers(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gettedCouriers))
	assert.Equal(t, freeCourier.ID(), gettedCouriers[0].ID())
}
