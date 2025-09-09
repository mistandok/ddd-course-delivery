package commands

import (
	"context"
	"errors"
	"testing"

	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrderHandler_Handle_SuccessfulOrderCreation(t *testing.T) {
	// Arrange
	mockOrderRepo := setupSuccessfulOrderRepo(t)
	mockUoW := setupSuccessfulUoW(t, mockOrderRepo)
	mockUoWFactory := setupUoWFactory(t, mockUoW)

	handler := NewCreateOrderHandler(mockUoWFactory)
	command := createValidCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

func TestCreateOrderHandler_Handle_InvalidCommand(t *testing.T) {
	// Arrange
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	handler := NewCreateOrderHandler(mockUoWFactory)
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
	mockUoW := setupSuccessfulUoW(t, mockOrderRepo)
	mockUoWFactory := setupUoWFactory(t, mockUoW)

	handler := NewCreateOrderHandler(mockUoWFactory)
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

	handler := NewCreateOrderHandler(mockUoWFactory)
	command := createValidCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

// Helper functions
func setupSuccessfulOrderRepo(t *testing.T) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
	return mockOrderRepo
}

func setupFailingOrderRepo(t *testing.T, expectedError error) *mocks.OrderRepo {
	mockOrderRepo := mocks.NewOrderRepo(t)
	mockOrderRepo.EXPECT().Add(mock.Anything, mock.Anything).Return(expectedError)
	return mockOrderRepo
}

func setupSuccessfulUoW(t *testing.T, orderRepo *mocks.OrderRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().OrderRepo().Return(orderRepo)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupFailingUoW(t *testing.T, expectedError error) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).Return(expectedError)
	return mockUoW
}

func setupUoWFactory(t *testing.T, uow *mocks.UnitOfWork) *mocks.UnitOfWorkFactory {
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockUoWFactory.EXPECT().NewUOW().Return(uow)
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
