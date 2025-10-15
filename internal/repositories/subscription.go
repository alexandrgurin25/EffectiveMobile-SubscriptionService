package repositories

import (
	"context"
	"fmt"
	"strings"
	"subscriptions/internal/entity"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, sub *entity.Subscription) (*entity.Subscription, error)
	GetById(ctx context.Context, id string) (*entity.Subscription, error)
	UpdateById(ctx context.Context, sub *entity.Subscription) error
	DeleteById(ctx context.Context, id string) error
	GetList(ctx context.Context, offset, limit int, userID, serviceName string) ([]entity.Subscription, error)
	CalculateSummary(ctx context.Context, userID, serviceName, startDate, endDate string) (int, error)
}

type subRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return &subRepository{pool: pool}
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
