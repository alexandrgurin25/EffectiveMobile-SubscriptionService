package services

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) GetById(ctx context.Context, id string) (*entity.Subscription, error) {

	subOut, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return subOut, nil
}
