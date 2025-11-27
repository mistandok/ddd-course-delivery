package create_order

import (
	"context"
	"errors"
	"testing"

	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrderHandler_Handle_SuccessfulOrderCreation(t *testing.T) {
	// Arrange
	mockOrderRepo := setupSuccessfulOrderRepo(t)
	mockGeoClient := setupSuccessfulGeoClient(t)
	mockUoW := setupSuccessfulUoW(t, mockOrderRepo)
	mockUoWFactory := setupUoWFactory(t, mockUoW)

	handler := NewCreateOrderHandler(mockUoWFactory, mockGeoClient)
	command := createValidCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

func TestCreateOrderHandler_Handle_InvalidCommand(t *testing.T) {
	// Arrange
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockGeoClient := mocks.NewGeoClient(t)
	handler := NewCreateOrderHandler(mockUoWFactory, mockGeoClient)
	command := createInvalidCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrCommandIsInvalid)
}

func TestCreateOrderHandler_Handle_OrderRepositoryError(t *testing.T) {
	// Arrange
	expectedError := errors.New("repository error")
	mockOrderRepo := setupFailingOrderRepo(t, expectedError)
	mockGeoClient := setupSuccessfulGeoClient(t)
	mockUoW := setupSuccessfulUoW(t, mockOrderRepo)
	mockUoWFactory := setupUoWFactory(t, mockUoW)

	handler := NewCreateOrderHandler(mockUoWFactory, mockGeoClient)
	command := createValidCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestCreateOrderHandler_Handle_UnitOfWorkDoError(t *testing.T) {
	// Arrange
	expectedError := errors.New("uow error")
	mockUoW := setupFailingUoW(t, expectedError)
	mockUoWFactory := setupUoWFactory(t, mockUoW)
	mockGeoClient := mocks.NewGeoClient(t)

	handler := NewCreateOrderHandler(mockUoWFactory, mockGeoClient)
	command := createValidCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestCreateOrderHandler_Handle_GeoClientError(t *testing.T) {
	// Arrange
	expectedError := errors.New("geo service error")
	mockGeoClient := setupFailingGeoClient(t, expectedError)
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoWFactory := setupUoWFactory(t, mockUoW)

	handler := NewCreateOrderHandler(mockUoWFactory, mockGeoClient)
	command := createValidCommand()

	mockUoW.On("Do", mock.Anything, mock.Anything).Return(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

// Helper functions
func setupSuccessfulOrderRepo(t *testing.T) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
	return mockOrderRepo
}

func setupFailingOrderRepo(t *testing.T, expectedError error) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.On("Add", mock.Anything, mock.Anything).Return(expectedError)
	return mockOrderRepo
}

func setupSuccessfulUoW(t *testing.T, orderRepo *mocks.OrderRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.On("OrderRepo").Return(orderRepo)
	mockUoW.On("Do", mock.Anything, mock.Anything).Return(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupFailingUoW(t *testing.T, expectedError error) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.On("Do", mock.Anything, mock.Anything).Return(expectedError)
	return mockUoW
}

func setupUoWFactory(t *testing.T, uow *mocks.UnitOfWork) *mocks.UnitOfWorkFactory {
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockUoWFactory.On("NewUOW").Return(uow)
	return mockUoWFactory
}

func createValidCommand() CreateOrderCommand {
	command, _ := NewCreateOrderCommand(uuid.New(), "test street", 10)
	return command
}

func createInvalidCommand() CreateOrderCommand {
	return CreateOrderCommand{
		orderID: uuid.New(),
		street:  "test street",
		volume:  10,
		isValid: false,
	}
}

func setupSuccessfulGeoClient(t *testing.T) *mocks.GeoClient {
	mockGeoClient := mocks.NewGeoClient(t)
	location, _ := shared_kernel.NewLocation(5, 5)
	mockGeoClient.On("GetGeolocation", mock.Anything).Return(location, nil)
	return mockGeoClient
}

func setupFailingGeoClient(t *testing.T, expectedError error) *mocks.GeoClient {
	mockGeoClient := mocks.NewGeoClient(t)
	mockGeoClient.On("GetGeolocation", mock.Anything).Return(shared_kernel.Location{}, expectedError)
	return mockGeoClient
}
