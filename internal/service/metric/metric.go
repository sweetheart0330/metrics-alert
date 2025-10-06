package metric

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/sweetheart0330/metrics-alert/internal/repository/interfaces"
	"go.uber.org/zap"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

var (
	ErrMetricNotFound    = errors.New("metric not found")
	ErrUnknownMetricType = errors.New("unknown metric type")
)

type Metric struct {
	repo interfaces.IRepository
	log  zap.SugaredLogger
}

func New(repo interfaces.IRepository, log zap.SugaredLogger) (*Metric, error) {
	metric := &Metric{
		repo: repo,
		log:  log,
	}

	return metric, nil
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
			m.log.Errorw("failed to get counter", "error", err)
			return models.Metrics{}, fmt.Errorf("counter value is nil")
		}
	case models.Gauge:
		if respMetric.Value == nil {
			m.log.Errorw("failed to get gauge", "error", err)
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

func (m *Metric) Ping(ctx context.Context) error {
	err := m.repo.Ping(ctx)
	if err != nil {
		m.log.Errorw("failed to ping metrics", "error", err)
		return fmt.Errorf("failed to ping: %w", err)
	}

	return nil
}
