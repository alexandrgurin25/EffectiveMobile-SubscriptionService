// sub_repository_test.go
package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"subscriptions/internal/entity"
	"subscriptions/internal/repositories/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSubRepository_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	userId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	sub := &entity.Subscription{
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    userId,
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	// parseDateToDB возвращает строки
	expectedStartDateStr := "2025-01-01"
	expectedEndDateStr := "2025-12-01"

	// Настраиваем ожидание вызова QueryRow
	mockDB.EXPECT().
		QueryRow(
			ctx,
			gomock.Any(), // SQL
			"Yandex Plus", 1500, userId, expectedStartDateStr, expectedEndDateStr,
		).
		Return(mockRow)

	// Для Scan нам нужны time.Time объекты
	expectedStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEndDate := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	// Настраиваем поведение Scan
	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			// Простое присваивание без проверок
			*(dest[0].(*string)) = "sub-123"            // id
			*(dest[1].(*string)) = "Yandex Plus"        // service_name
			*(dest[2].(*int)) = 1500                    // price
			*(dest[3].(*string)) = userId               // user_id
			*(dest[4].(*time.Time)) = expectedStartDate // start_date
			*(dest[5].(*sql.NullTime)) = sql.NullTime{  // end_date
				Time:  expectedEndDate,
				Valid: true,
			}
			return nil
		})

	// Вызываем тестируемый метод
	result, err := repo.Create(ctx, sub)

	// Проверяем результат
	require.NoError(t, err)
	assert.Equal(t, "sub-123", result.Id)
	assert.Equal(t, "Yandex Plus", result.Name)
	assert.Equal(t, 1500, result.Price)
	assert.Equal(t, userId, result.UserId)
	assert.Equal(t, "01-2025", result.StartDate)
	assert.Equal(t, "12-2025", result.EndDate)
}
