package repositories

import (
	"context"
	"database/sql"
	"errors"
	"subscriptions/internal/entity"
	"subscriptions/internal/repositories/mocks"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubRepository_UpdateById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "sub-123",
		Name:      "Yandex Plus Premium",
		Price:     2000,
		UserId:    "user-123",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	// Ожидаемые преобразованные даты
	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	// Настраиваем ожидание вызова Exec
	mockDB.EXPECT().
		Exec(
			ctx,
			gomock.Any(), // SQL
			"sub-123", "Yandex Plus Premium", 2000, "user-123", expectedStartDateStr, expectedEndDateStr,
		).
		Return(pgconn.NewCommandTag("UPDATE 1"), nil)

	// Вызываем тестируемый метод
	err := repo.UpdateById(ctx, sub)

	// Проверяем что ошибки нет
	assert.NoError(t, err)
}

func TestSubRepository_UpdateById_WithoutEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "sub-123",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "user-123",
		StartDate: "01-2025",
		EndDate:   "", // Пустая end date
	}

	// Ожидаемые преобразованные даты
	expectedStartDateStr := "2025-01-01"

	// Настраиваем ожидание вызова Exec с nil для end_date
	mockDB.EXPECT().
		Exec(
			ctx,
			gomock.Any(), // SQL
			"sub-123", "Yandex Plus", 1500, "user-123", expectedStartDateStr, nil,
		).
		Return(pgconn.NewCommandTag("UPDATE 1"), nil)

	// Вызываем тестируемый метод
	err := repo.UpdateById(ctx, sub)

	// Проверяем что ошибки нет
	assert.NoError(t, err)
}

func TestSubRepository_UpdateById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "non-existent-id",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "user-123",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	// Настраиваем ожидание вызова Exec с возвратом sql.ErrNoRows
	mockDB.EXPECT().
		Exec(
			ctx,
			gomock.Any(),
			"non-existent-id", "Yandex Plus", 1500, "user-123", expectedStartDateStr, expectedEndDateStr,
		).
		Return(pgconn.NewCommandTag("UPDATE 0"), sql.ErrNoRows)

	// Вызываем тестируемый метод
	err := repo.UpdateById(ctx, sub)

	// Проверяем что вернулась ошибка sql.ErrNoRows
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestSubRepository_UpdateById_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "sub-123",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "user-123",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	// Настраиваем ожидание вызова Exec с возвратом ошибки БД
	mockDB.EXPECT().
		Exec(
			ctx,
			gomock.Any(),
			"sub-123", "Yandex Plus", 1500, "user-123", expectedStartDateStr, expectedEndDateStr,
		).
		Return(pgconn.NewCommandTag(""), assert.AnError)

	// Вызываем тестируемый метод
	err := repo.UpdateById(ctx, sub)

	// Проверяем что вернулась ошибка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to UPDATE subscription")
}

func TestSubRepository_UpdateById_InvalidStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "sub-123",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "user-123",
		StartDate: "invalid-date", // Невалидная дата
		EndDate:   "12-2025",
	}

	// Вызываем тестируемый метод
	err := repo.UpdateById(ctx, sub)

	// Проверяем что вернулась ошибка валидации start date
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start date")
}

func TestSubRepository_UpdateById_InvalidEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "sub-123",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "user-123",
		StartDate: "01-2025",
		EndDate:   "invalid-date", // Невалидная дата
	}

	// Вызываем тестируемый метод
	err := repo.UpdateById(ctx, sub)

	// Проверяем что вернулась ошибка валидации end date
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid end date")
}
