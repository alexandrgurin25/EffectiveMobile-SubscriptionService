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

func TestGetById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	subId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	expectedSub := &entity.Subscription{
		Id:        subId,
		Name:      "Yandex Plus",
		Price:     1500,
		UserId:    "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate: "01-2025",
		EndDate:   "12-2025",
	}

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetById(ctx, subId).
		Return(expectedSub, nil).Times(1)

	service := New(mockRepo)

	sub, err := service.GetById(ctx, subId)

	require.NoError(t, err)          
	assert.Equal(t, expectedSub, sub) 
}

func TestGetById_Fail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	subId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetById(ctx, subId).
		Return(nil, fmt.Errorf("subscription not found")).Times(1)

	service := New(mockRepo)

	sub, err := service.GetById(ctx, subId)

	assert.Error(t, err) 
	assert.Equal(t, "subscription not found", err.Error())
	assert.Nil(t, sub) 
}
