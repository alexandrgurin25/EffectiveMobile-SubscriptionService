package services

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) GetById(ctx context.Context, id string) (*entity.Subscription, error) {

	return s.repo.GetById(ctx, id)
}
