package services

import (
	"context"
	"subscriptions/internal/entity"
	"subscriptions/internal/repositories"
)

type Service interface {
	Create(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error)
	GetById(ctx context.Context, id string) (*entity.Subscription, error)
	UpdateById(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error)
	DeleteById(ctx context.Context, id string) error
	GetList(ctx context.Context, id string) ([]entity.Subscription, error)
}

type subService struct {
	repo repositories.Repository
}

func New(repo repositories.Repository) Service {
	return &subService{repo: repo}
}
