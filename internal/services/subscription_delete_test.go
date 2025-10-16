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

func TestDeleteById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	subId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetById(ctx, subId).
		Return(&entity.Subscription{Id: subId}, nil).Times(1)

	mockRepo.EXPECT().DeleteById(ctx, subId).
		Return(nil).Times(1)

	service := New(mockRepo)

	err := service.DeleteById(ctx, subId)

	
	require.NoError(t, err)
}

func TestDeleteById_Fail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	subId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetById(ctx, subId).
		Return(nil, fmt.Errorf("subscription not found")).Times(1)

	service := New(mockRepo)

	err := service.DeleteById(ctx, subId)

	require.Error(t, err)
	assert.Equal(t, "subscription not found", err.Error())
}

func TestDeleteById_Fail_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	subId := "60601fee-2bf1-4721-ae6f-7636e79a0cba"

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetById(ctx, subId).
		Return(&entity.Subscription{Id: subId}, nil).Times(1)

	mockRepo.EXPECT().DeleteById(ctx, subId).
		Return(fmt.Errorf("deletion error")).Times(1)

	service := New(mockRepo)

	err := service.DeleteById(ctx, subId)

	require.Error(t, err) 
	assert.Equal(t, "deletion error", err.Error())
}
