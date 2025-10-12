package client

import models "github.com/sweetheart0330/metrics-alert/internal/model"

//go:generate mockgen -source=./interface.go -destination=./../mocks/mock_client.go -package=mocks
type IClient interface {
	SendGaugeMetric(m models.Metrics) error
	SendCounterMetric(m models.Metrics) error
	SendMetricsBatch(metrics []models.Metrics) error
}
