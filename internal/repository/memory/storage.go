package memory

import (
	"errors"
	"sync"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
)

type MemStorage struct {
	metrics map[string]models.Metrics
	mu      sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]models.Metrics),
	}
}

func (ms *MemStorage) UpdateGaugeMetric(metric models.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.metrics[metric.ID] = metric

	return nil
}

func (ms *MemStorage) UpdateCounterMetric(metric models.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

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

func (ms *MemStorage) UpdateAllMetrics(metrics []models.Metrics) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, val := range metrics {
		ms.metrics[val.ID] = val
	}
}

func (ms *MemStorage) GetMetric(metricID string) (models.Metrics, error) {
	m, ok := ms.metrics[metricID]
	if !ok {
		return models.Metrics{}, metric.ErrMetricNotFound
	}

	return m, nil
}
func (ms *MemStorage) GetAllMetrics() ([]models.Metrics, error) {
	mList := make([]models.Metrics, 0, len(ms.metrics))
	for _, m := range ms.metrics {
		mList = append(mList, m)
	}

	return mList, nil
}
