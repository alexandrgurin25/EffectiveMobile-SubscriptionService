package services

import (
	"context"
)

func (s *subService) DeleteById(ctx context.Context, id string) error {

	_, err := s.repo.GetById(ctx, id)

	if err != nil {
		return err
	}

	err = s.repo.DeleteById(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
