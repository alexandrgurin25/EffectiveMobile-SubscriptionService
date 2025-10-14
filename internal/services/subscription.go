package service

import (
	"context"
	"subscriptions/internal/entity"
	"subscriptions/internal/repositories"
)

type Service interface {
	Create(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error)
}

type subService struct {
	repo repositories.Repository
}

func New(repo repositories.Repository) Service {
	return &subService{repo: repo}
}
