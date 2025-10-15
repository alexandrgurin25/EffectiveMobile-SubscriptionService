package services

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) GetList(ctx context.Context, id string) ([]entity.Subscription, error) {

	subsOut, err := s.repo.GetListByUserID(ctx, id)

	if err != nil {
		return nil, err
	}

	return subsOut, nil
}
