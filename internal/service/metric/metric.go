package metric

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/observer"
	"github.com/sweetheart0330/metrics-alert/internal/repository/interfaces"
	"go.uber.org/zap"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

var (
	ErrMetricNotFound    = errors.New("metric not found")
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrConnRepo          = errors.New("connection repository error")
)

type Metric struct {
	repo interfaces.IRepository
	log  zap.SugaredLogger
	obs  observer.Publisher
}

func New(repo interfaces.IRepository, log zap.SugaredLogger, observer observer.Publisher) (*Metric, error) {
	metric := &Metric{
		repo: repo,
		log:  log,
		obs:  observer,
	}

	return metric, nil
}

func (m *Metric) UpdateMetric(ctx context.Context, metric models.Metrics) error {
	switch metric.MType {
	case models.Counter:
		err := m.repo.UpdateCounterMetric(ctx, metric)
		if err != nil {
			return fmt.Errorf("failed to update counter metric: %w", err)
		}

		m.sendAuditFromSingleMetric(ctx, metric)
		return nil
	case models.Gauge:
		err := m.repo.UpdateGaugeMetric(ctx, metric)
		if err != nil {
			return fmt.Errorf("failed to update gauge metric: %w", err)
		}

		m.sendAuditFromSingleMetric(ctx, metric)
		return nil
	}

	return fmt.Errorf("%w: %s", ErrUnknownMetricType, metric.MType)
}

func (m *Metric) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	err := m.repo.UpdateMetrics(ctx, metrics)
	if err != nil {
		m.log.Errorw("failed to update metrics", "error", err)
		return fmt.Errorf("failed to update metrics: %w", err)
	}

	ev := m.formMetrics(ctx, metrics)
	m.obs.NotifyObservers(ctx, ev)

	return nil
}

func (m *Metric) GetMetric(ctx context.Context, metricID string) (models.Metrics, error) {
	respMetric, err := m.repo.GetMetric(ctx, metricID)
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

func (m *Metric) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	metrics, err := m.repo.GetAllMetrics(ctx)
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

func (m *Metric) sendAuditFromSingleMetric(ctx context.Context, metric models.Metrics) {
	var metrics []models.Metrics
	metrics = append(metrics, metric)
	ev := m.formMetrics(ctx, metrics)
	m.obs.NotifyObservers(ctx, ev)
}

func (m *Metric) formMetrics(ctx context.Context, metrics []models.Metrics) models.AuditEvent {
	metricNames := make([]string, len(metrics))

	for i, m := range metrics {
		metricNames[i] = m.ID
	}

	return models.AuditEvent{
		TS:        time.Now().Unix(),
		Metrics:   metricNames,
		IPAddress: ctx.Value(models.CtxClientIP).(string),
	}
}
