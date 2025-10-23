package postgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	servMetric "github.com/sweetheart0330/metrics-alert/internal/service/metric"
)

const (
	checkMetricTable = `SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables 
            WHERE table_schema = 'public' 
              AND table_name = 'metrics'
        );`
	GetMetricsQuery = `SELECT metric_id, metric_type, delta_value, gauge_value
						FROM metrics
						WHERE metric_id = $1
						ORDER BY timestamp DESC
						LIMIT 1;
						`
	UpdateMetricsCounterQuery = `INSERT INTO metrics (metric_id, metric_type, delta_value) 
                   VALUES ($1, $2, $3)
                   ON CONFLICT (metric_id) DO UPDATE
                   SET
                       metric_type = EXCLUDED.metric_type,
                       delta_value = metrics.delta_value + EXCLUDED.delta_value;`
	UpdateMetricsGaugeQuery = `INSERT INTO metrics (metric_id, metric_type, gauge_value) 
							VALUES ($1, $2, $3)
							ON CONFLICT (metric_id) DO UPDATE
							SET
						    	metric_type = EXCLUDED.metric_type,
						    	gauge_value  = EXCLUDED.gauge_value;
						`
	GetAllMetrics = `SELECT metric_id, metric_type, delta_value, gauge_value FROM metrics`
)

func (db *Database) UpdateCounterMetric(ctx context.Context, metric models.Metrics) error {
	_, err := db.pg.Exec(ctx, UpdateMetricsCounterQuery, metric.ID, metric.MType, metric.Delta)
	if err != nil {
		return handleError(fmt.Errorf("failed to create/update counter: %w", err))
	}

	return nil
}

func (db *Database) UpdateGaugeMetric(ctx context.Context, metric models.Metrics) error {
	_, err := db.pg.Exec(ctx, UpdateMetricsGaugeQuery, metric.ID, metric.MType, metric.Value)
	if err != nil {
		return handleError(fmt.Errorf("failed to create/update gauge: %w", err))
	}

	return nil
}

func (db *Database) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	tx, err := db.pg.Begin(ctx)
	if err != nil {
		return handleError(fmt.Errorf("failed to start transaction, err: %w", err))
	}

	for _, metric := range metrics {
		switch metric.MType {
		case models.Gauge:
			_, err = db.pg.Exec(ctx, UpdateMetricsGaugeQuery, metric.ID, metric.MType, metric.Value)
			if err != nil {
				return handleError(fmt.Errorf("failed to create/update gauge metric: %w", err))
			}
		case models.Counter:
			_, err = db.pg.Exec(ctx, UpdateMetricsCounterQuery, metric.ID, metric.MType, metric.Delta)
			if err != nil {
				return handleError(fmt.Errorf("failed to create/update counter metric: %w", err))
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handleError(fmt.Errorf("failed to commit transaction, err: %w", err))
	}

	return nil
}

func (db *Database) GetMetric(ctx context.Context, metricID string) (models.Metrics, error) {
	m := models.Metrics{}
	err := db.pg.QueryRow(ctx, GetMetricsQuery, metricID).Scan(
		&m.ID,
		&m.MType,
		&m.Delta,
		&m.Value,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metrics{}, servMetric.ErrMetricNotFound
		}

		return models.Metrics{}, fmt.Errorf("failed to send query: %w", err)
	}

	return m, nil
}

func (db *Database) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	rows, err := db.pg.Query(ctx, GetAllMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to send query: %w", err)
	}
	defer rows.Close()

	var metrics []models.Metrics
	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(
			&m.ID,
			&m.MType,
			&m.Delta,
			&m.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return metrics, nil
}

func handleError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Проверяем на Connection Exception
		if pgerrcode.IsConnectionException(pgErr.Code) {
			log.Printf("Ошибка подключения: %s (код: %s)", pgErr.Message, pgErr.Code)
			// Здесь можно добавить логику переподключения
			return fmt.Errorf("%w, err: %w", servMetric.ErrConnRepo, err)
		}
	}

	return err
}
