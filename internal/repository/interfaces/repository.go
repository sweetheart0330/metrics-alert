package interfaces

import (
	"context"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

//go:generate mockgen -source=./repository.go -destination=./../mocks/mock_repo.go -package=mocks
type IRepository interface {
	UpdateCounterMetric(metric models.Metrics) error
	UpdateGaugeMetric(metric models.Metrics) error
	UpdateAllMetrics(metrics []models.Metrics)
	GetMetric(metricID string) (models.Metrics, error)
	GetAllMetrics() ([]models.Metrics, error)
	Ping(ctx context.Context) error
}
