package contracts

import models "github.com/sweetheart0330/metrics-alert/internal/model"

type MetricService interface {
	UpdateGaugeMetric(metrics models.Metrics) error
	UpdateCounterMetric(metrics models.Metrics) error
}
