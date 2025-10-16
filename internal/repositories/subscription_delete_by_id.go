package repositories

import (
	"context"
	"fmt"
)

func (r *subRepository) DeleteById(ctx context.Context, id string) error {

	_, err := r.db.Exec(
		ctx,
		`DELETE FROM subscriptions 
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to DELETE subscription: %v", err)
	}

	return nil
}
