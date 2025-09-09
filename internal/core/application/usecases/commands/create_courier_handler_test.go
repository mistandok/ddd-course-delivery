package commands

import (
	"context"
	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/errs"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCourierHandler_Handle_SuccessfulCourierCreation(t *testing.T) {
	// Arrange
	mockCourierRepo := setupSuccessfulCourierRepo(t)
	mockUoW := setupSuccessfulUoWForCourier(t, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForCourier(t, mockUoW)
	
	handler := NewCreateCourierHandler(mockUoWFactory)
	command := createValidCourierCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

func TestCreateCourierHandler_Handle_InvalidCommand(t *testing.T) {
	// Arrange
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	handler := NewCreateCourierHandler(mockUoWFactory)
	command := createInvalidCourierCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrCommandIsInvalid)
}

func TestCreateCourierHandler_Handle_CourierRepositoryError(t *testing.T) {
	// Arrange
	expectedError := errors.New("repository error")
	mockCourierRepo := setupFailingCourierRepo(t, expectedError)
	mockUoW := setupSuccessfulUoWForCourier(t, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForCourier(t, mockUoW)
	
	handler := NewCreateCourierHandler(mockUoWFactory)
	command := createValidCourierCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestCreateCourierHandler_Handle_UnitOfWorkDoError(t *testing.T) {
	// Arrange
	expectedError := errors.New("uow error")
	mockUoW := setupFailingUoWForCourier(t, expectedError)
	mockUoWFactory := setupUoWFactoryForCourier(t, mockUoW)
	
	handler := NewCreateCourierHandler(mockUoWFactory)
	command := createValidCourierCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

// Helper functions
func setupSuccessfulCourierRepo(t *testing.T) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
	return mockCourierRepo
}

func setupFailingCourierRepo(t *testing.T, expectedError error) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Add(mock.Anything, mock.Anything).Return(expectedError)
	return mockCourierRepo
}

func setupSuccessfulUoWForCourier(t *testing.T, courierRepo *mocks.CourierRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().CourierRepo().Return(courierRepo)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupFailingUoWForCourier(t *testing.T, expectedError error) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).Return(expectedError)
	return mockUoW
}

func setupUoWFactoryForCourier(t *testing.T, uow *mocks.UnitOfWork) *mocks.UnitOfWorkFactory {
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockUoWFactory.EXPECT().NewUOW().Return(uow)
	return mockUoWFactory
}

func createValidCourierCommand() CreateCourierCommand {
	command, _ := NewCreateCourierCommand("Test Courier", 50)
	return command
}

func createInvalidCourierCommand() CreateCourierCommand {
	return CreateCourierCommand{
		name:    "Test Courier",
		speed:   50,
		isValid: false,
	}
}