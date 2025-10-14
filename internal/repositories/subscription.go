package repositories

import (
	"context"
	"subscriptions/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error)
}

type subRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return &subRepository{pool: pool}
}
