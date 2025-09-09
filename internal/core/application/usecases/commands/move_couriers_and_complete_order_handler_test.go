package commands

import (
	"context"
	"errors"
	"testing"

	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/errs"

	modelCourier "delivery/internal/core/domain/model/courier"
	modelOrder "delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoveCouriersAndFinishOrderHandler_Handle_SuccessfulMovementAndCompletion(t *testing.T) {
	// Arrange
	order := newValidAssignedOrder(t)
	courier := newValidCourierForMovement(t, order)

	mockOrderRepo := setupSuccessfulOrderRepoWithAssignedOrders(t, []*modelOrder.Order{order})
	mockCourierRepo := setupSuccessfulCourierRepoForMovement(t, courier)
	mockUoW := setupSuccessfulUoWForMovement(t, mockOrderRepo, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_InvalidCommand(t *testing.T) {
	// Arrange
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createInvalidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrCommandIsInvalid)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_GetAssignedOrdersError(t *testing.T) {
	// Arrange
	expectedError := errors.New("get assigned orders error")
	mockOrderRepo := setupFailingOrderRepoForGetAssigned(t, expectedError)
	mockUoW := setupUoWWithOrderRepo(t, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_GetCourierError(t *testing.T) {
	// Arrange
	order := newValidAssignedOrder(t)
	expectedError := errors.New("get courier error")

	mockOrderRepo := setupSuccessfulOrderRepoWithAssignedOrders(t, []*modelOrder.Order{order})
	mockCourierRepo := setupFailingCourierRepoForGetInMovement(t, *order.CourierID(), expectedError)
	mockUoW := setupUoWWithBothRepos(t, mockOrderRepo, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_UpdateCourierError(t *testing.T) {
	// Arrange
	order := newValidAssignedOrder(t)
	courier := newValidCourierForMovement(t, order)
	expectedError := errors.New("update courier error")

	mockOrderRepo := setupSuccessfulOrderRepoWithAssignedOrders(t, []*modelOrder.Order{order})
	mockCourierRepo := setupCourierRepoWithGetSuccessUpdateFailure(t, courier, expectedError)
	mockUoW := setupUoWWithBothRepos(t, mockOrderRepo, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_UpdateOrderError(t *testing.T) {
	// Arrange
	order := newValidAssignedOrder(t)
	courier := newValidCourierForMovement(t, order)
	expectedError := errors.New("update order error")

	mockOrderRepo := setupOrderRepoWithGetSuccessUpdateFailure(t, []*modelOrder.Order{order}, expectedError)
	mockCourierRepo := setupSuccessfulCourierRepoForMovement(t, courier)
	mockUoW := setupUoWWithBothRepos(t, mockOrderRepo, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_UnitOfWorkDoError(t *testing.T) {
	// Arrange
	expectedError := errors.New("uow error")
	mockUoW := setupFailingUoWForMovement(t, expectedError)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestMoveCouriersAndFinishOrderHandler_Handle_NoAssignedOrders(t *testing.T) {
	// Arrange
	mockOrderRepo := setupSuccessfulOrderRepoWithAssignedOrders(t, []*modelOrder.Order{})
	mockUoW := setupUoWWithOrderRepo(t, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForMovement(t, mockUoW)

	handler := NewMoveCouriersAndCompleteOrderHandler(mockUoWFactory)
	command := createValidMoveCouriersCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

// Helper functions
func newValidAssignedOrder(t *testing.T) *modelOrder.Order {
	t.Helper()
	orderLocation, _ := shared_kernel.NewLocation(5, 5)
	order, _ := modelOrder.NewOrder(uuid.New(), orderLocation, 10)
	courierID := uuid.New()
	_ = order.Assign(courierID)
	return order
}

func newValidCourierForMovement(t *testing.T, order *modelOrder.Order) *modelCourier.Courier {
	t.Helper()
	courierLocation, _ := shared_kernel.NewLocation(1, 1)
	courier, _ := modelCourier.NewCourier("Test Courier", 10, courierLocation)
	_ = courier.TakeOrder(order)
	return courier
}

func setupSuccessfulOrderRepoWithAssignedOrders(t *testing.T, orders []*modelOrder.Order) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().GetAllInAssignedStatus(mock.Anything).Return(orders, nil)
	mockOrderRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Maybe()
	return mockOrderRepo
}

func setupFailingOrderRepoForGetAssigned(t *testing.T, expectedError error) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().GetAllInAssignedStatus(mock.Anything).Return(nil, expectedError)
	return mockOrderRepo
}

func setupOrderRepoWithGetSuccessUpdateFailure(t *testing.T, orders []*modelOrder.Order, expectedError error) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().GetAllInAssignedStatus(mock.Anything).Return(orders, nil)
	mockOrderRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(expectedError)
	return mockOrderRepo
}

func setupSuccessfulCourierRepoForMovement(t *testing.T, courier *modelCourier.Courier) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Get(mock.Anything, mock.Anything).Return(courier, nil)
	mockCourierRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	return mockCourierRepo
}

func setupFailingCourierRepoForGetInMovement(t *testing.T, courierID uuid.UUID, expectedError error) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Get(mock.Anything, courierID).Return(nil, expectedError)
	return mockCourierRepo
}

func setupCourierRepoWithGetSuccessUpdateFailure(t *testing.T, courier *modelCourier.Courier, expectedError error) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Get(mock.Anything, mock.Anything).Return(courier, nil)
	mockCourierRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(expectedError)
	return mockCourierRepo
}

func setupSuccessfulUoWForMovement(t *testing.T, orderRepo *mocks.OrderRepo, courierRepo *mocks.CourierRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().OrderRepo().Return(orderRepo)
	mockUoW.EXPECT().CourierRepo().Return(courierRepo)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupUoWWithOrderRepo(t *testing.T, orderRepo *mocks.OrderRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().OrderRepo().Return(orderRepo)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupUoWWithBothRepos(t *testing.T, orderRepo *mocks.OrderRepo, courierRepo *mocks.CourierRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().OrderRepo().Return(orderRepo)
	mockUoW.EXPECT().CourierRepo().Return(courierRepo)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupFailingUoWForMovement(t *testing.T, expectedError error) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).Return(expectedError)
	return mockUoW
}

func setupUoWFactoryForMovement(t *testing.T, uow *mocks.UnitOfWork) *mocks.UnitOfWorkFactory {
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockUoWFactory.EXPECT().NewUOW().Return(uow)
	return mockUoWFactory
}

func createValidMoveCouriersCommand() MoveCouriersAndFinishOrderCommand {
	return NewMoveCouriersAndFinishOrderCommand()
}

func createInvalidMoveCouriersCommand() MoveCouriersAndFinishOrderCommand {
	return MoveCouriersAndFinishOrderCommand{
		isValid: false,
	}
}
