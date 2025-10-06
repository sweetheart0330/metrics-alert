package contracts

import (
	"context"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

//go:generate mockgen -source=./metric.go -destination=./../../mocks/mock_metric.go -package=mocks
type MetricService interface {
	UpdateMetric(ctx context.Context, metrics models.Metrics) error
	GetMetric(ctx context.Context, metricID string) (models.Metrics, error)
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
	Ping(ctx context.Context) error
}
