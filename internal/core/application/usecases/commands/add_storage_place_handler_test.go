package commands

import (
	"context"
	"errors"
	"testing"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/shared_kernel"
	"delivery/internal/core/ports/mocks"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddStoragePlaceHandler_Handle_SuccessfulStoragePlaceAddition(t *testing.T) {
	// Arrange
	courierID := uuid.New()
	testCourier := newValidCourier(t)

	mockCourierRepo := setupSuccessfulCourierRepoForGet(t, courierID, testCourier)
	mockCourierRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	mockUoW := setupSuccessfulUoWForAddStoragePlace(t, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForAddStoragePlace(t, mockUoW)

	handler := NewAddStoragePlaceHandler(mockUoWFactory)
	command := createValidAddStoragePlaceCommand(courierID)

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.NoError(t, err)
}

func TestAddStoragePlaceHandler_Handle_InvalidCommand(t *testing.T) {
	// Arrange
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	handler := NewAddStoragePlaceHandler(mockUoWFactory)
	command := createInvalidAddStoragePlaceCommand()

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrCommandIsInvalid)
}

func TestAddStoragePlaceHandler_Handle_CourierNotFound(t *testing.T) {
	// Arrange
	courierID := uuid.New()
	expectedError := errors.New("courier not found")

	mockCourierRepo := setupFailingCourierRepoForGet(t, courierID, expectedError)
	mockUoW := setupSuccessfulUoWForAddStoragePlace(t, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForAddStoragePlace(t, mockUoW)

	handler := NewAddStoragePlaceHandler(mockUoWFactory)
	command := createValidAddStoragePlaceCommand(courierID)

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAddStoragePlaceHandler_Handle_CourierRepositoryUpdateError(t *testing.T) {
	// Arrange
	courierID := uuid.New()
	testCourier := newValidCourier(t)
	expectedError := errors.New("update error")

	mockCourierRepo := setupSuccessfulCourierRepoForGet(t, courierID, testCourier)
	mockCourierRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(expectedError)

	mockUoW := setupSuccessfulUoWForAddStoragePlace(t, mockCourierRepo)
	mockUoWFactory := setupUoWFactoryForAddStoragePlace(t, mockUoW)

	handler := NewAddStoragePlaceHandler(mockUoWFactory)
	command := createValidAddStoragePlaceCommand(courierID)

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestAddStoragePlaceHandler_Handle_UnitOfWorkDoError(t *testing.T) {
	// Arrange
	courierID := uuid.New()
	expectedError := errors.New("uow error")

	mockUoW := setupFailingUoWForAddStoragePlace(t, expectedError)
	mockUoWFactory := setupUoWFactoryForAddStoragePlace(t, mockUoW)

	handler := NewAddStoragePlaceHandler(mockUoWFactory)
	command := createValidAddStoragePlaceCommand(courierID)

	// Act
	err := handler.Handle(context.Background(), command)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

// Helper functions
func setupSuccessfulCourierRepoForGet(t *testing.T, courierID uuid.UUID, testCourier *courier.Courier) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Get(mock.Anything, courierID).Return(testCourier, nil)
	return mockCourierRepo
}

func setupFailingCourierRepoForGet(t *testing.T, courierID uuid.UUID, expectedError error) *mocks.CourierRepo {
	mockCourierRepo := mocks.NewCourierRepo(t)
	mockCourierRepo.EXPECT().Get(mock.Anything, courierID).Return(nil, expectedError)
	return mockCourierRepo
}

func setupSuccessfulUoWForAddStoragePlace(t *testing.T, courierRepo *mocks.CourierRepo) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().CourierRepo().Return(courierRepo)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	return mockUoW
}

func setupFailingUoWForAddStoragePlace(t *testing.T, expectedError error) *mocks.UnitOfWork {
	mockUoW := mocks.NewUnitOfWork(t)
	mockUoW.EXPECT().Do(mock.Anything, mock.Anything).Return(expectedError)
	return mockUoW
}

func setupUoWFactoryForAddStoragePlace(t *testing.T, uow *mocks.UnitOfWork) *mocks.UnitOfWorkFactory {
	mockUoWFactory := mocks.NewUnitOfWorkFactory(t)
	mockUoWFactory.EXPECT().NewUOW().Return(uow)
	return mockUoWFactory
}

func createValidAddStoragePlaceCommand(courierID uuid.UUID) AddStoragePlaceCommand {
	command, _ := NewAddStoragePlaceCommand(courierID, "Test Storage", 100)
	return command
}

func createInvalidAddStoragePlaceCommand() AddStoragePlaceCommand {
	return AddStoragePlaceCommand{
		courierID:   uuid.New(),
		name:        "Test Storage",
		totalVolume: 100,
		isValid:     false,
	}
}

func newValidCourier(t *testing.T) *courier.Courier {
	t.Helper()

	location, err := shared_kernel.NewRandomLocation()
	if err != nil {
		t.Fatalf("failed to create random location: %v", err)
	}

	testCourier, err := courier.NewCourier("Test Courier", 50, location)
	if err != nil {
		t.Fatalf("failed to create courier: %v", err)
	}

	return testCourier
}
