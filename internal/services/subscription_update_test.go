package services

import (
	"context"
	"fmt"
	"subscriptions/internal/entity"
	"subscriptions/internal/repositories/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().UpdateById(ctx, sub).Return(nil).Times(1)
	mockRepo.EXPECT().GetById(ctx, sub.Id).Return(sub, nil).Times(1)

	service := New(mockRepo)
	updatedSub, err := service.UpdateById(ctx, sub)

	require.NoError(t, err)
	assert.Equal(t, sub, updatedSub)
}

func TestUpdateById_Fail_UpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().UpdateById(ctx, sub).Return(fmt.Errorf("update error")).Times(1)

	service := New(mockRepo)
	updatedSub, err := service.UpdateById(ctx, sub)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, updatedSub)
}

func TestUpdateById_Fail_GetByIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	sub := &entity.Subscription{
		Id:        "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().UpdateById(ctx, sub).Return(nil).Times(1)
	mockRepo.EXPECT().GetById(ctx, sub.Id).Return(nil, fmt.Errorf("subscription not found")).Times(1)

	service := New(mockRepo)
	updatedSub, err := service.UpdateById(ctx, sub)

	assert.Error(t, err)
	assert.Equal(t, "subscription not found", err.Error())
	assert.Nil(t, updatedSub)
}
