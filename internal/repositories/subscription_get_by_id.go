package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscriptions/internal/entity"
	"time"
)

func (r *subRepository) GetById(ctx context.Context, id string) (*entity.Subscription, error) {

	var sub entity.Subscription

	var startDateDB time.Time
	var endDateDB sql.NullTime

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions 
		WHERE id = $1`,
		id,
	).Scan(&sub.Id, &sub.Name, &sub.Price, &sub.UserId, &startDateDB, &endDateDB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to GET subscription: %v", err)
	}

	sub.StartDate = formatTimeToMMYYYY(startDateDB)

	if endDateDB.Valid {
		sub.EndDate = formatTimeToMMYYYY(endDateDB.Time)
	} else {
		sub.EndDate = ""
	}

	return &sub, nil
}
