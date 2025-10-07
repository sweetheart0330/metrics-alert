package interfaces

import (
	"context"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

//go:generate mockgen -source=./repository.go -destination=./../../mocks/mock_repo.go -package=mocks
type IRepository interface {
	UpdateCounterMetric(ctx context.Context, metric models.Metrics) error
	UpdateGaugeMetric(ctx context.Context, metric models.Metrics) error
	GetMetric(ctx context.Context, metricID string) (models.Metrics, error)
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
	Ping(ctx context.Context) error
}
