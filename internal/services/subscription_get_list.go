package services

import (
	"context"
	"subscriptions/internal/entity"
)

func (s *subService) GetList(ctx context.Context, page, limit int,
	userID, serviceName string) ([]entity.Subscription, bool, error) {

	offset := (page - 1) * limit

	subs, err := s.repo.GetList(ctx, offset, limit+1, userID, serviceName)
	if err != nil {
		return nil, false, err
	}

	hasNext := false
	if len(subs) > limit {
		hasNext = true
		subs = subs[:limit]
	}

	return subs, hasNext, nil
}
