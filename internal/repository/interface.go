package repository

import (
	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../mocks/mock_repo.go
type IRepository interface {
	UpdateCounterMetric(metric models.Metrics) error
	UpdateGaugeMetric(metric models.Metrics) error
	GetMetric(metricID string) (models.Metrics, error)
	GetAllMetrics() ([]models.Metrics, error)
}

type FileSaver interface {
	WriteMetrics(metrics []models.Metrics) error
	UploadMetrics() ([]models.Metrics, error)
}
