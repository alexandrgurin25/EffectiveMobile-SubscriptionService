package services

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) Create(ctx context.Context, subIn *entity.Subscription) (*entity.Subscription, error) {

	return s.repo.Create(ctx, subIn)

}
