package memory

import (
	"fmt"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
)

type MemStorage struct {
	metrics map[string]models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: map[string]models.Metrics{},
	}
}

func (ms *MemStorage) UpdateGaugeMetric(metric models.Metrics) error {
	ms.metrics[metric.ID] = metric

	fmt.Printf("current gauge ID: %s, value: %f\n", metric.ID, *metric.Value)
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(metric models.Metrics) error {
	val, ok := ms.metrics[metric.ID]
	if !ok {
		ms.metrics[metric.ID] = metric
		fmt.Printf("current gauge ID: %s, value: %d\n", metric.ID, *metric.Delta)
		return nil
	}

	*val.Delta = *val.Delta + *metric.Delta
	ms.metrics[metric.ID] = val

	fmt.Println("current counter value: ", ms.metrics[metric.ID])
	return nil
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
