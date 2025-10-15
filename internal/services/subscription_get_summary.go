package services

import "context"

func (s *subService) GetSummary(ctx context.Context, userId string, serviceName string,
	startDate string, endDate string) (int, error) {

	return s.repo.CalculateSummary(ctx, userId, serviceName, startDate, endDate)
}
