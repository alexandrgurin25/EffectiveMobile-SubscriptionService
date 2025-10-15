package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscriptions/internal/entity"
)

func (r *subRepository) UpdateById(ctx context.Context, subIn *entity.Subscription) error {

	startDateForDB, err := parseDateToDB(subIn.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %v", err)
	}

	var endDateForDB interface{}
	if subIn.EndDate == "" {
		endDateForDB = nil
	} else {
		parsedEndDate, err := parseDateToDB(subIn.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date: %v", err)
		}
		endDateForDB = parsedEndDate
	}

	_, err = r.pool.Exec(
		ctx,
		`Update subscriptions 
		SET service_name = $2, price = $3, user_id = $4, start_date = $5, end_date = $6
		WHERE id = $1`,
		subIn.Id,
		subIn.Name,
		subIn.Price,
		subIn.UserId,
		startDateForDB,
		endDateForDB,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return fmt.Errorf("failed to UPDATE subscription: %v", err)
	}

	return nil
}
