package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"subscriptions/internal/entity"
	"time"
)

func (r *subRepository) Create(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error) {
	var price int
	var id, name, userId string
	var startDateDB time.Time
	var endDateDB sql.NullTime

	startDateForDB, err := parseDateToDB(sub.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	var endDateForDB interface{}
	if sub.EndDate == "" {
		endDateForDB = nil
	} else {
		parsedEndDate, err := parseDateToDB(sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %v", err)
		}
		endDateForDB = parsedEndDate
	}

	err = r.pool.QueryRow(
		ctx,
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date`,
		sub.Name,
		sub.Price,
		sub.UserId,
		startDateForDB,
		endDateForDB,
	).Scan(&id, &name, &price, &userId, &startDateDB, &endDateDB)

	if err != nil {
		return nil, fmt.Errorf("failed to CREATE subscription: %v", err)
	}

	startDateFormatted := formatTimeToMMYYYY(startDateDB)

	var endDateFormatted string
	if endDateDB.Valid {
		endDateFormatted = formatTimeToMMYYYY(endDateDB.Time)
	} else {
		endDateFormatted = ""
	}

	return &entity.Subscription{
		Id:        id,
		Name:      name,
		Price:     price,
		UserId:    userId,
		StartDate: startDateFormatted,
		EndDate:   endDateFormatted,
	}, nil
}

func parseDateToDB(dateStr string) (string, error) {
	if dateStr == "" {
		return "", nil
	}

	if len(dateStr) == 7 && strings.Contains(dateStr, "-") {
		parts := strings.Split(dateStr, "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid date format: %s", dateStr)
		}

		month := parts[0]
		if month < "01" || month > "12" {
			return "", fmt.Errorf("invalid month: %s", month)
		}

		return fmt.Sprintf("%s-%s-01", parts[1], parts[0]), nil
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %s", dateStr)
	}

	return dateStr, nil
}

func formatTimeToMMYYYY(t time.Time) string {
	return t.Format("01-2006")
}
