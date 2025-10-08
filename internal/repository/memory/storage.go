package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/repository/interfaces"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
	"go.uber.org/zap"
)

type MemStorage struct {
	metrics   map[string]models.Metrics
	mu        sync.Mutex
	log       *zap.SugaredLogger
	fileSaver interfaces.FileSaver

	storeInterval uint
}

func NewMemStorage(ctx context.Context, fileSaver interfaces.FileSaver, log *zap.SugaredLogger, restore bool, storeInterval uint) *MemStorage {
	mem := &MemStorage{
		metrics:       make(map[string]models.Metrics),
		fileSaver:     fileSaver,
		log:           log,
		storeInterval: storeInterval,
	}

	if restore {
		metrics, err := mem.fileSaver.UploadMetrics()
		if err != nil {
			mem.log.Warnw("failed to upload metrics", "error", err)
		} else {
			mem.UpdateAllMetrics(ctx, metrics)
			mem.log.Debug("metrics are successfully restored")
		}
	}

	if storeInterval > 0 {
		go mem.saveInPeriod(ctx)
	}

	return mem
}

func (ms *MemStorage) saveInPeriod(ctx context.Context) {
	t := time.NewTicker(time.Duration(ms.storeInterval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			err := ms.saveToFile()
			if err != nil {
				ms.log.Errorw("failed to save to file", "error", err)
			}
			ms.log.Info("saved metrics to file in period")
		}
	}
}

func (ms *MemStorage) saveToFile() error {
	metrics, err := ms.GetAllMetrics(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get metrics, err: %w", err)
	}
	err = ms.fileSaver.WriteMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to write metrics, err: %w", err)
	}

	return nil
}

func (ms *MemStorage) UpdateGaugeMetric(_ context.Context, metric models.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.metrics[metric.ID] = metric

	return nil
}

func (ms *MemStorage) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	for _, m := range metrics {
		switch m.MType {
		case models.Counter:
			err := ms.UpdateGaugeMetric(ctx, m)
			if err != nil {
				return fmt.Errorf("failed to update counter metric, err: %w", err)
			}
		case models.Gauge:
			err := ms.UpdateGaugeMetric(ctx, m)
			if err != nil {
				return fmt.Errorf("failed to update gauge metric, err: %w", err)
			}
		}
	}

	return nil
}

func (ms *MemStorage) UpdateCounterMetric(_ context.Context, metric models.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.storeInterval == 0 {
		err := ms.saveToFile()
		if err != nil {
			return fmt.Errorf("failed to save to file, err: %w", err)
		}
	}

	if metric.Delta == nil {
		return errors.New("metrics counter delta is nil")
	}

	val, ok := ms.metrics[metric.ID]
	if !ok {
		ms.metrics[metric.ID] = metric
		return nil
	}

	*val.Delta += *metric.Delta
	ms.metrics[metric.ID] = val

	return nil
}

func (ms *MemStorage) UpdateAllMetrics(_ context.Context, metrics []models.Metrics) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, val := range metrics {
		ms.metrics[val.ID] = val
	}
}

func (ms *MemStorage) GetMetric(_ context.Context, metricID string) (models.Metrics, error) {
	m, ok := ms.metrics[metricID]
	if !ok {
		return models.Metrics{}, metric.ErrMetricNotFound
	}

	return m, nil
}
func (ms *MemStorage) GetAllMetrics(_ context.Context) ([]models.Metrics, error) {
	mList := make([]models.Metrics, 0, len(ms.metrics))
	for _, m := range ms.metrics {
		mList = append(mList, m)
	}

	return mList, nil
}

func (ms *MemStorage) Ping(_ context.Context) error {
	// в данном пакете нечего пинговать
	return nil
}
