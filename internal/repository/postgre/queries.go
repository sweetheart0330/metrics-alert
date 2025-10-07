package postgre

import (
	"context"
	"fmt"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
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
	UpdateMetricsQuery = `INSERT INTO metrics (metric_id, metric_type, delta_value, gauge_value) 
							VALUES ($1, $2, $3, $4)
							ON CONFLICT (metric_id) DO UPDATE
							SET
						    	metric_type = EXCLUDED.metric_type,
						    	delta_value  = EXCLUDED.delta_value,
						    	gauge_value  = EXCLUDED.gauge_value;
						`
	GetAllMetrics = `SELECT metric_id, metric_type, delta_value, gauge_value FROM metrics`
)

func (db *Database) UpdateCounterMetric(ctx context.Context, metric models.Metrics) error {
	_, err := db.pg.Exec(ctx, UpdateMetricsQuery, metric.ID, metric.MType, metric.Delta, nil)
	if err != nil {
		return fmt.Errorf("failed to create/update counter: %w", err)
	}

	return nil
}

func (db *Database) UpdateGaugeMetric(ctx context.Context, metric models.Metrics) error {
	_, err := db.pg.Exec(ctx, UpdateMetricsQuery, metric.ID, metric.MType, nil, metric.Value)
	if err != nil {
		return fmt.Errorf("failed to create/update gauge: %w", err)
	}

	return nil
}

func (db *Database) GetMetric(ctx context.Context, metricID string) (models.Metrics, error) {
	m := models.Metrics{}
	err := db.pg.QueryRow(ctx, GetMetricsQuery, metricID).Scan(
		m.ID,
		m.MType,
		m.Delta,
		m.Value,
	)
	if err != nil {
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
			m.ID,
			m.MType,
			m.Delta,
			m.Value)
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
