package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"subscriptions/internal/entity"
	"time"
)

func (r *subRepository) GetList(ctx context.Context, offset, limit int, userID, serviceName string) ([]entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, serviceName)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription
	for rows.Next() {
		var startDateDB time.Time
		var endDateDB sql.NullTime
		var s entity.Subscription
		err := rows.Scan(&s.Id, &s.Name, &s.Price, &s.UserId, &startDateDB, &endDateDB)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		s.StartDate = formatTimeToMMYYYY(startDateDB)
		if endDateDB.Valid {
			s.EndDate = formatTimeToMMYYYY(endDateDB.Time)
		}

		subs = append(subs, s)
	}

	if len(subs) == 0 {
		return subs, sql.ErrNoRows
	}
	
	return subs, nil
}
