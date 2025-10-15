package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"subscriptions/internal/entity"
	"time"
)

func (r *subRepository) GetListByUserID(ctx context.Context, userID string) ([]entity.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions
        WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {
		var sub entity.Subscription
		var startDateDB time.Time
		var endDateDB sql.NullTime

		if err := rows.Scan(
			&sub.Id,
			&sub.Name,
			&sub.Price,
			&sub.UserId,
			&startDateDB,
			&endDateDB,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		sub.StartDate = formatTimeToMMYYYY(startDateDB)
		if endDateDB.Valid {
			sub.EndDate = formatTimeToMMYYYY(endDateDB.Time)
		}

		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if len(subs) == 0 {
		return subs, sql.ErrNoRows
	}

	return subs, nil
}
