package repositories

import (
	"context"
	"testing"

	"subscriptions/internal/repositories/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubRepository_DeleteById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	subscriptionID := "sub-123"

	commandTag := pgconn.NewCommandTag("DELETE 1")

	mockDB.EXPECT().
		Exec(
			ctx,
			gomock.Any(),
			subscriptionID,
		).
		Return(commandTag, nil)

	err := repo.DeleteById(ctx, subscriptionID)
	assert.NoError(t, err)
}

func TestSubRepository_DeleteById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	repo := &subRepository{db: mockDB}

	ctx := context.Background()
	subscriptionID := "non-existent-id"

	commandTag := pgconn.NewCommandTag("DELETE 0")

	mockDB.EXPECT().
		Exec(
			ctx,
			gomock.Any(),
			subscriptionID,
		).
		Return(commandTag, nil)

	err := repo.DeleteById(ctx, subscriptionID)
	assert.NoError(t, err)
}

func TestSubRepository_DeleteById_DBError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDB := mocks.NewMockDB(ctrl)
    repo := &subRepository{db: mockDB}

    ctx := context.Background()
    subscriptionID := "sub-123"

    emptyCommandTag := pgconn.NewCommandTag("")

    mockDB.EXPECT().
        Exec(
            ctx,
            gomock.Any(),
            subscriptionID,
        ).
        Return(emptyCommandTag, assert.AnError)

    err := repo.DeleteById(ctx, subscriptionID)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to DELETE subscription")
}
