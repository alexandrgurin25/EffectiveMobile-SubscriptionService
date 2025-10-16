package services

import (
	"context"
	"fmt"
	"subscriptions/internal/repositories/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetSummary_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "12-2025"

	expectedSummary := 3200

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().CalculateSummary(ctx, userId, serviceName, startDate, endDate).
		Return(expectedSummary, nil).Times(1)

	service := New(mockRepo)

	summary, err := service.GetSummary(ctx, userId, serviceName, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, expectedSummary, summary)
}

func TestGetSummary_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	serviceName := "Yandex Plus"
	startDate := "01-2025"
	endDate := "12-2025"

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().CalculateSummary(ctx, userId, serviceName, startDate, endDate).
		Return(0, fmt.Errorf("internal server error")).Times(1)

	service := New(mockRepo)

	summary, err := service.GetSummary(ctx, userId, serviceName, startDate, endDate)

	assert.Error(t, err)
	assert.Equal(t, "internal server error", err.Error())
	assert.Equal(t, 0, summary)
}
