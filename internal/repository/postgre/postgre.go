package postgre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sweetheart0330/metrics-alert/internal/repository/interfaces"
)

type Database struct {
	pg *pgxpool.Pool
}

func NewDatabase(ctx context.Context, connStr string) (interfaces.IRepository, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{pg: pool}, nil
}

func (db *Database) Close() {
	db.pg.Close()
}

func (db *Database) Ping(ctx context.Context) error {
	err := db.pg.Ping(ctx)
	if err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	return nil
}
