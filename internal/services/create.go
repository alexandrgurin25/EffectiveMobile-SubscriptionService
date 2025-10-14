package service

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) Create(ctx context.Context, subIn *entity.Subscription) (*entity.Subscription, error) {
	subOut, err := s.repo.Create(ctx, subIn)

	if err != nil {
		return nil, err
	}

	return subOut, nil

}
