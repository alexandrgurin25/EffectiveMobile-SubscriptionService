package repositories

import (
	"context"
	"subscriptions/internal/repositories/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubRepository_GetList_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	offset := 0
	limit := 10

	mockDB.EXPECT().
		Query(
			ctx,
			gomock.Any(),
			limit, offset,
		).
		Return(nil, assert.AnError)

	result, err := repo.GetList(ctx, offset, limit, "", "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
