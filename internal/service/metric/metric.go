package metric

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/repository"
)

var (
	ErrMetricNotFound    = errors.New("metric not found")
	ErrUnknownMetricType = errors.New("unknown metric type")
)

type Metric struct {
	repo repository.IRepository
}

func New(repo repository.IRepository) *Metric {
	return &Metric{repo: repo}
}

func (m *Metric) UpdateMetric(metrics models.Metrics) error {
	switch metrics.MType {
	case models.Counter:
		return m.repo.UpdateCounterMetric(metrics)
	case models.Gauge:
		return m.repo.UpdateGaugeMetric(metrics)
	}

	return fmt.Errorf("%w: %s", ErrUnknownMetricType, metrics.MType)
}

func (m *Metric) GetMetric(metric string) (models.Metrics, error) {
	respMetric, err := m.repo.GetMetric(metric)
	if err != nil {
		return models.Metrics{}, fmt.Errorf("failed to get metric: %w", err)
	}

	switch respMetric.MType {
	case models.Counter:
		if respMetric.Delta == nil {
			fmt.Println("failed to get counter, err: ", err)
			return models.Metrics{}, fmt.Errorf("counter value is nil")
		}
	case models.Gauge:
		if respMetric.Value == nil {
			fmt.Println("failed to get gauge, err: ", err)
			return models.Metrics{}, fmt.Errorf("gauge value is nil")
		}
	}

	return respMetric, nil
}

func (m *Metric) GetAllMetrics() ([]models.Metrics, error) {
	metrics, err := m.repo.GetAllMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	// сортировка по ID
	if len(metrics) != 0 {
		slices.SortFunc(metrics, func(a, b models.Metrics) int {
			return strings.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID))
		})
	}

	return metrics, nil
}
