// database.go
package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)


//go:generate mockgen -destination=mocks/row_mock.go -package=mocks github.com/jackc/pgx/v5 Row
type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type PgxPoolAdapter struct {
	pool *pgxpool.Pool
}

func NewPgxPoolAdapter(pool *pgxpool.Pool) *PgxPoolAdapter {
	return &PgxPoolAdapter{pool: pool}
}

func (a *PgxPoolAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return a.pool.QueryRow(ctx, sql, args...)
}

func (a *PgxPoolAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return a.pool.Exec(ctx, sql, args...)
}

func (a *PgxPoolAdapter) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return a.pool.Query(ctx, sql, args...)
}
