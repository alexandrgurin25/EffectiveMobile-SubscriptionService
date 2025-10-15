package repositories

import (
	"context"
	"fmt"
)

func (r *subRepository) CalculateSummary(ctx context.Context, userID, serviceName, startDate, endDate string) (int, error) {
	var totalCost int

	startDateForDB, err := parseDateToDB(startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date: %v", err)
	}

	query := `
		SELECT COALESCE(SUM(price), 0) 
		FROM subscriptions 
		WHERE user_id = $1 AND service_name = $2 
		AND start_date >= $3
	`

	args := []interface{}{userID, serviceName, startDateForDB}

	if endDate != "" {
		query += " AND (end_date <= $4 OR end_date IS NULL)"
		endDateForDB, err := parseDateToDB(endDate)
		if err != nil {
			return 0, fmt.Errorf("invalid start date: %v", err)
		}
		args = append(args, endDateForDB)
	} else {
		query += " AND (end_date IS NULL OR end_date >= $4)"
		args = append(args, startDateForDB)
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate summary: %w", err)
	}

	return totalCost, nil
}
