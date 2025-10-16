package repositories

import (
	"context"
	"subscriptions/internal/repositories/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSubRepository_CalculateSummary_SuccessWithEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userID := "user-123"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "12-2025"

	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(),
			userID, serviceName, expectedStartDateStr, expectedEndDateStr,
		).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*int)) = 4500 
			return nil
		})

	totalCost, err := repo.CalculateSummary(ctx, userID, serviceName, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, 4500, totalCost)
}

func TestSubRepository_CalculateSummary_SuccessWithoutEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userID := "user-123"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "" 

	expectedStartDateStr := "2025-01-01"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(),
			userID, serviceName, expectedStartDateStr, expectedStartDateStr,
		).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*int)) = 3000 
			return nil
		})

	totalCost, err := repo.CalculateSummary(ctx, userID, serviceName, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, 3000, totalCost)
}

func TestSubRepository_CalculateSummary_ZeroResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userID := "user-123"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "12-2025"

	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(),
			userID, serviceName, expectedStartDateStr, expectedEndDateStr,
		).
		Return(mockRow)

	// COALESCE вернет 0 если нет подписок
	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*int)) = 0
			return nil
		})

	totalCost, err := repo.CalculateSummary(ctx, userID, serviceName, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, 0, totalCost)
}

func TestSubRepository_CalculateSummary_InvalidStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userID := "user-123"
	serviceName := "Yandex Plus"
	startDate := "invalid-date" // Невалидная дата
	endDate := "12-2025"

	totalCost, err := repo.CalculateSummary(ctx, userID, serviceName, startDate, endDate)

	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start date")
	assert.Equal(t, 0, totalCost)
}

func TestSubRepository_CalculateSummary_InvalidEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userID := "user-123"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "invalid-date" 

	totalCost, err := repo.CalculateSummary(ctx, userID, serviceName, startDate, endDate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start date")
	assert.Equal(t, 0, totalCost)
}

func TestSubRepository_CalculateSummary_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userID := "user-123"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "12-2025"

	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(),
			userID, serviceName, expectedStartDateStr, expectedEndDateStr,
		).
		Return(mockRow)

	// Эмулируем ошибку БД
	mockRow.EXPECT().
		Scan(gomock.Any()).
		Return(assert.AnError)

	totalCost, err := repo.CalculateSummary(ctx, userID, serviceName, startDate, endDate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to calculate summary")
	assert.Equal(t, 0, totalCost)
}
