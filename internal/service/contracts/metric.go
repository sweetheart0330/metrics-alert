package contracts

import (
	"context"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

type MetricService interface {
	UpdateMetric(metrics models.Metrics) error
	GetMetric(metric string) (models.Metrics, error)
	GetAllMetrics() ([]models.Metrics, error)
	Ping(ctx context.Context) error
}
