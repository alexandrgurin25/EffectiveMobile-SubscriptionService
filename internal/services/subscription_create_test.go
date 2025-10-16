package services

import (
	"context"
	"fmt"

	"subscriptions/internal/entity"
	"subscriptions/internal/repositories/mocks"

	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateSubscription_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	userId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	subIn := &entity.Subscription{
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    userId,
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	subOut := &entity.Subscription{
		Id:        "d6d273fa-486e-4d74-94e0-94dd9b95a1d8",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    userId,
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().Create(ctx, subIn).
		Return(subOut, nil)

	service := New(mockRepo)

	sub, err := service.Create(ctx, subIn)
	require.NoError(t, err)
	require.Equal(t, sub, subOut)
}

func TestCreateSubscription_Fail_Validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	subIn := &entity.Subscription{
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    userId,
		StartDate: "01-2025",
		EndDate:   "invalid", 
	}

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().Create(ctx, subIn).
		Return(nil, fmt.Errorf("invalid date range")).Times(1)

	service := New(mockRepo)

	
	sub, err := service.Create(ctx, subIn)

	require.Error(t, err) 
	require.Equal(t, "invalid date range", err.Error())
	require.Nil(t, sub) 
}
