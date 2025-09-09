package assign_order

import (
	"context"
	"errors"
	"testing"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssignedOrderHandler_Handle_SuccessfulOrderAssignment(t *testing.T) {
	// Arrange
	testOrder := newValidOrder(t)
	testCouriers := newValidCouriers(t)
	selectedCourier := testCouriers[0]

	mockCourierRepo := setupSuccessfulCourierRepoForAssignment(t, testCouriers)
	mockCourierRepo.EXPECT().Update(mock.Anything, selectedCourier).Return(nil)

	mockOrderRepo := setupSuccessfulOrderRepoForAssignment(t, testOrder)
	mockOrderRepo.EXPECT().Update(mock.Anything, testOrder).Return(nil)

	mockOrderDispatcher := setupSuccessfulOrderDispatcher(t, testOrder, testCouriers, selectedCourier)

	mockUoW := setupSuccessfulUoWForAssignment(t, mockCourierRepo, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

func TestAssignedOrderHandler_Handle_InvalidCommand(t *testing.T) {
	// Arrange
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockOrderDispatcher := mocks.NewOrderDispatcher(t)
	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createInvalidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrCommandIsInvalid)
}

func TestAssignedOrderHandler_Handle_GetAllFreeCouriersError(t *testing.T) {
	// Arrange
	expectedError := errors.New("failed to get free couriers")
	mockCourierRepo := setupFailingCourierRepoForGetAllFree(t, expectedError)
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderDispatcher := mocks.NewOrderDispatcher(t)

	mockUoW := setupSuccessfulUoWForAssignment(t, mockCourierRepo, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAssignedOrderHandler_Handle_GetFirstInCreatedStatusError(t *testing.T) {
	// Arrange
	testCouriers := newValidCouriers(t)
	expectedError := errors.New("failed to get order in created status")

	mockCourierRepo := setupSuccessfulCourierRepoForAssignment(t, testCouriers)
	mockOrderRepo := setupFailingOrderRepoForGetFirst(t, expectedError)
	mockOrderDispatcher := mocks.NewOrderDispatcher(t)

	mockUoW := setupSuccessfulUoWForAssignment(t, mockCourierRepo, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAssignedOrderHandler_Handle_OrderDispatcherError(t *testing.T) {
	// Arrange
	testOrder := newValidOrder(t)
	testCouriers := newValidCouriers(t)
	expectedError := errors.New("dispatch failed")

	mockCourierRepo := setupSuccessfulCourierRepoForAssignment(t, testCouriers)
	mockOrderRepo := setupSuccessfulOrderRepoForAssignment(t, testOrder)
	mockOrderDispatcher := setupFailingOrderDispatcher(t, expectedError)

	mockUoW := setupSuccessfulUoWForAssignment(t, mockCourierRepo, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAssignedOrderHandler_Handle_OrderRepositoryUpdateError(t *testing.T) {
	// Arrange
	testOrder := newValidOrder(t)
	testCouriers := newValidCouriers(t)
	selectedCourier := testCouriers[0]
	expectedError := errors.New("order update failed")

	mockCourierRepo := setupSuccessfulCourierRepoForAssignment(t, testCouriers)
	mockOrderRepo := setupSuccessfulOrderRepoForAssignment(t, testOrder)
	mockOrderRepo.EXPECT().Update(mock.Anything, testOrder).Return(expectedError)

	mockOrderDispatcher := setupSuccessfulOrderDispatcher(t, testOrder, testCouriers, selectedCourier)

	mockUoW := setupSuccessfulUoWForAssignment(t, mockCourierRepo, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAssignedOrderHandler_Handle_CourierRepositoryUpdateError(t *testing.T) {
	// Arrange
	testOrder := newValidOrder(t)
	testCouriers := newValidCouriers(t)
	selectedCourier := testCouriers[0]
	expectedError := errors.New("courier update failed")

	mockCourierRepo := setupSuccessfulCourierRepoForAssignment(t, testCouriers)
	mockCourierRepo.EXPECT().Update(mock.Anything, selectedCourier).Return(expectedError)

	mockOrderRepo := setupSuccessfulOrderRepoForAssignment(t, testOrder)
	mockOrderRepo.EXPECT().Update(mock.Anything, testOrder).Return(nil)

	mockOrderDispatcher := setupSuccessfulOrderDispatcher(t, testOrder, testCouriers, selectedCourier)

	mockUoW := setupSuccessfulUoWForAssignment(t, mockCourierRepo, mockOrderRepo)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAssignedOrderHandler_Handle_UnitOfWorkDoError(t *testing.T) {
	// Arrange
	expectedError := errors.New("uow error")
	mockUoW := setupFailingUoWForAssignment(t, expectedError)
	mockUoWFactory := setupUoWFactoryForAssignment(t, mockUoW)
	mockOrderDispatcher := mocks.NewOrderDispatcher(t)

	handler := NewAssignedOrderHandler(mockUoWFactory, mockOrderDispatcher)
	command := createValidAssignedOrderCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

// Helper functions
func setupSuccessfulCourierRepoForAssignment(t *testing.T, testCouriers []*courier.Courier) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().GetAllFreeCouriers(mock.Anything).Return(testCouriers, nil)
	return mockCourierRepo
}

func setupFailingCourierRepoForGetAllFree(t *testing.T, expectedError error) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().GetAllFreeCouriers(mock.Anything).Return(nil, expectedError)
	return mockCourierRepo
}

func setupSuccessfulOrderRepoForAssignment(t *testing.T, testOrder *order.Order) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().GetFirstInCreatedStatus(mock.Anything).Return(testOrder, nil)
	return mockOrderRepo
}

func setupFailingOrderRepoForGetFirst(t *testing.T, expectedError error) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().GetFirstInCreatedStatus(mock.Anything).Return(nil, expectedError)
	return mockOrderRepo
}

func setupSuccessfulOrderDispatcher(t *testing.T, testOrder *order.Order, testCouriers []*courier.Courier, selectedCourier *courier.Courier) *mocks.OrderDispatcher {
	mockOrderDispatcher := mocks.NewOrderDispatcher(t)
	mockOrderDispatcher.EXPECT().Dispatch(testOrder, testCouriers).Return(selectedCourier, nil)
	return mockOrderDispatcher
}

func setupFailingOrderDispatcher(t *testing.T, expectedError error) *mocks.OrderDispatcher {
	mockOrderDispatcher := mocks.NewOrderDispatcher(t)
	mockOrderDispatcher.EXPECT().Dispatch(mock.Anything, mock.Anything).Return(nil, expectedError)
	return mockOrderDispatcher
}

func setupSuccessfulUoWForAssignment(t *testing.T, courierRepo *mocks.CourierRepo, orderRepo *mocks.OrderRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().CourierRepo().Return(courierRepo).Maybe()
	mockUoW.EXPECT().OrderRepo().Return(orderRepo).Maybe()
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupFailingUoWForAssignment(t *testing.T, expectedError error) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).Return(expectedError)
	return mockUoW
}

func setupUoWFactoryForAssignment(t *testing.T, uow *mocks.UnitOfWork) *mocks.UnitOfWorkFactory {
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockUoWFactory.EXPECT().NewUOW().Return(uow)
	return mockUoWFactory
}

func createValidAssignedOrderCommand() AssignedOrderCommand {
	return NewAssignedOrderCommand()
}

func createInvalidAssignedOrderCommand() AssignedOrderCommand {
	return AssignedOrderCommand{
		isValid: false,
	}
}

func newValidOrder(t *testing.T) *order.Order {
	t.Helper()

	location, err := shared_kernel.NewLocation(5, 8)
	if err != nil {
		t.Fatalf("failed to create location: %v", err)
	}
	testOrder, err := order.NewOrder(uuid.New(), location, 15)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	return testOrder
}

func newValidCouriers(t *testing.T) []*courier.Courier {
	t.Helper()

	location1, err := shared_kernel.NewLocation(3, 7)
	if err != nil {
		t.Fatalf("failed to create location1: %v", err)
	}
	courier1, err := courier.NewCourier("Test Courier 1", 50, location1)
	if err != nil {
		t.Fatalf("failed to create courier 1: %v", err)
	}

	location2, err := shared_kernel.NewLocation(8, 4)
	if err != nil {
		t.Fatalf("failed to create location2: %v", err)
	}
	courier2, err := courier.NewCourier("Test Courier 2", 40, location2)
	if err != nil {
		t.Fatalf("failed to create courier 2: %v", err)
	}

	return []*courier.Courier{courier1, courier2}
}
