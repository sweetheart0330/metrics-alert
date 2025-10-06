package metric

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/repository"
)

var (
	ErrMetricNotFound    = errors.New("metric not found")
	ErrUnknownMetricType = errors.New("unknown metric type")
)

type Metric struct {
	repo          repository.IRepository
	fileStorage   repository.FileSaver
	log           zap.SugaredLogger
	storeInterval uint
	restore       bool
}

func New(ctx context.Context, repo repository.IRepository, fileStorage repository.FileSaver, storeInterval uint, restore bool, log zap.SugaredLogger) (*Metric, error) {
	metric := &Metric{
		repo:          repo,
		fileStorage:   fileStorage,
		storeInterval: storeInterval,
		log:           log,
	}

	if restore {
		metrics, err := metric.fileStorage.UploadMetrics()
		if err != nil {
			log.Warnw("failed to upload metrics", "error", err)
		} else {
			repo.UpdateAllMetrics(metrics)
			log.Debug("metrics are successfully restored")
		}
	}

	if storeInterval > 0 {
		go metric.saveInPeriod(ctx)
	}

	return metric, nil
}

func (m *Metric) saveInPeriod(ctx context.Context) {
	t := time.NewTicker(time.Duration(m.storeInterval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			err := m.saveToFile()
			if err != nil {
				m.log.Errorw("failed to save to file", "error", err)
			}
			m.log.Info("saved metrics to file in period")
		}
	}
}

func (m *Metric) UpdateMetric(metrics models.Metrics) error {
	if m.storeInterval == 0 {
		err := m.saveToFile()
		if err != nil {
			return fmt.Errorf("failed to save to file, err: %w", err)
		}
	}

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

func (m *Metric) saveToFile() error {
	metrics, err := m.repo.GetAllMetrics()
	if err != nil {
		return fmt.Errorf("failed to get metrics, err: %w", err)
	}
	err = m.fileStorage.WriteMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to write metrics, err: %w", err)
	}

	return nil
}
