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

func TestGetList_Success_WithNextPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	serviceName := "Yandex Plus"
	page := 1
	limit := 2

	expectedSubs := []entity.Subscription{
		{Id: "1", Name: "Yandex Plus", Price: 1500, UserId: userID},
		{Id: "2", Name: "Yandex Music", Price: 700, UserId: userID},
		{Id: "3", Name: "Yandex Video", Price: 500, UserId: userID},
	}

	mockRepo := mocks.NewMockRepository(ctrl)


	mockRepo.EXPECT().GetList(ctx, 0, limit+1, userID, serviceName).
		Return(expectedSubs, nil).Times(1)

	service := New(mockRepo)

	subs, hasNext, err := service.GetList(ctx, page, limit, userID, serviceName)

	require.NoError(t, err)
	assert.Equal(t, 2, len(subs))
	assert.True(t, hasNext)
}

func TestGetList_Success_WithoutNextPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	serviceName := "Yandex Plus"
	page := 1
	limit := 2

	expectedSubs := []entity.Subscription{
		{Id: "1", Name: "Yandex Plus", Price: 1500, UserId: userID},
		{Id: "2", Name: "Yandex Music", Price: 700, UserId: userID},
	}

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetList(ctx, 0, limit+1, userID, serviceName).
		Return(expectedSubs, nil).Times(1)

	service := New(mockRepo)

	subs, hasNext, err := service.GetList(ctx, page, limit, userID, serviceName)

	require.NoError(t, err)      
	assert.Equal(t, 2, len(subs))
	assert.False(t, hasNext)      
}

func TestGetList_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	serviceName := "Yandex Plus"
	page := 1
	limit := 2

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().GetList(ctx, 0, limit+1, userID, serviceName).
		Return(nil, fmt.Errorf("internal server error")).Times(1)

	service := New(mockRepo)

	subs, hasNext, err := service.GetList(ctx, page, limit, userID, serviceName)

	assert.Error(t, err)
	assert.Equal(t, "internal server error", err.Error())
	assert.Nil(t, subs)
	assert.False(t, hasNext)
}
