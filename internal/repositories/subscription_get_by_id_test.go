package repositories

import (
	"context"
	"database/sql"
	"errors"
	"subscriptions/internal/repositories/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSubRepository_GetById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	subscriptionID := "sub-123"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			`SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions 
		WHERE id = $1`,
			subscriptionID,
		).
		Return(mockRow)

	expectedStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEndDate := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*string)) = "sub-123"            // id
			*(dest[1].(*string)) = "Yandex Plus"        // service_name
			*(dest[2].(*int)) = 1500                    // price
			*(dest[3].(*string)) = "user-123"           // user_id
			*(dest[4].(*time.Time)) = expectedStartDate // start_date
			*(dest[5].(*sql.NullTime)) = sql.NullTime{  // end_date
				Time:  expectedEndDate,
				Valid: true,
			}
			return nil
		})

	result, err := repo.GetById(ctx, subscriptionID)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sub-123", result.Id)
	assert.Equal(t, "Yandex Plus", result.Name)
	assert.Equal(t, 1500, result.Price)
	assert.Equal(t, "user-123", result.UserId)
	assert.Equal(t, "01-2025", result.StartDate)
	assert.Equal(t, "12-2025", result.EndDate)
}

func TestSubRepository_GetById_WithoutEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	subscriptionID := "sub-123"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(), 
			subscriptionID,
		).
		Return(mockRow)

	expectedStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*string)) = "sub-123"            // id
			*(dest[1].(*string)) = "Yandex Plus"        // service_name
			*(dest[2].(*int)) = 1500                    // price
			*(dest[3].(*string)) = "user-123"           // user_id
			*(dest[4].(*time.Time)) = expectedStartDate // start_date
			*(dest[5].(*sql.NullTime)) = sql.NullTime{  // end_date
				Valid: false,
			}
			return nil
		})

	result, err := repo.GetById(ctx, subscriptionID)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "sub-123", result.Id)
	assert.Equal(t, "Yandex Plus", result.Name)
	assert.Equal(t, 1500, result.Price)
	assert.Equal(t, "user-123", result.UserId)
	assert.Equal(t, "01-2025", result.StartDate)
	assert.Equal(t, "", result.EndDate) 
}

func TestSubRepository_GetById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	subscriptionID := "non-existent-id"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(),
			subscriptionID,
		).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any()).
		Return(sql.ErrNoRows)

	result, err := repo.GetById(ctx, subscriptionID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	assert.Nil(t, result)
}

func TestSubRepository_GetById_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	subscriptionID := "sub-123"

	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(),
			subscriptionID,
		).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any()).
		Return(assert.AnError)

	result, err := repo.GetById(ctx, subscriptionID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to GET subscription")
}
