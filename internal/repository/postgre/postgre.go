package postgre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
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

func (db *Database) UpdateCounterMetric(metric models.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (db *Database) UpdateGaugeMetric(metric models.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (db *Database) UpdateAllMetrics(metrics []models.Metrics) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) GetMetric(metricID string) (models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) GetAllMetrics() ([]models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) Ping(ctx context.Context) error {
	err := db.pg.Ping(ctx)
	if err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	return nil
}
