package services

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) UpdateById(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error) {

	err := s.repo.Update(ctx, sub)

	if err != nil {
		return nil, err
	}

	subOut, err := s.repo.GetById(ctx, sub.Id)
	if err != nil {
		return nil, err
	}

	return subOut, nil
}
