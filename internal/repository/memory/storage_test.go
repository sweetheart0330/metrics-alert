package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

func Test_NewMemStorage(t *testing.T) {
	metricMap := make(map[string]models.Metrics)

	memStrg := NewMemStorage()

	assert.Equal(t, metricMap, memStrg.metrics)
}

func Test_UpdateGaugeMetric(t *testing.T) {
	memStrg := MemStorage{
		metrics: make(map[string]models.Metrics),
	}

	val := 12.5
	metric := models.Metrics{
		ID:    "test-name",
		Value: &val,
	}

	_ = memStrg.UpdateGaugeMetric(metric)

	assert.Equal(t, val, *memStrg.metrics[metric.ID].Value)
}

func Test_UpdateCounterMetric(t *testing.T) {
	memStrg := MemStorage{
		metrics: make(map[string]models.Metrics),
	}

	val := int64(12)
	metric := models.Metrics{
		ID:    "test-name",
		Delta: &val,
	}

	_ = memStrg.UpdateCounterMetric(metric)

	assert.Equal(t, val, *memStrg.metrics[metric.ID].Delta)

	//sec call
	val2 := int64(2)
	metric.Delta = &val2
	_ = memStrg.UpdateCounterMetric(metric)

	assert.Equal(t, int64(14), *memStrg.metrics[metric.ID].Delta)
}
