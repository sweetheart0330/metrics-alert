package metric

import (
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/repository"
)

type Metric struct {
	repo repository.IRepository
}

func New(repo repository.IRepository) *Metric {
	return &Metric{repo: repo}
}

func (m *Metric) UpdateGaugeMetric(metrics models.Metrics) error {
	return m.repo.UpdateGaugeMetric(metrics.ID, *metrics.Value)
}

func (m *Metric) UpdateCounterMetric(metrics models.Metrics) error {
	return m.repo.UpdateCounterMetric(metrics.ID, *metrics.Delta)
}
