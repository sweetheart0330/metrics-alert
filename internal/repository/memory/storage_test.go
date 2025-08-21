package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMemStorage(t *testing.T) {
	gaugeMap := make(map[string]float64)
	counterMap := make(map[string]int64)

	memStrg := NewMemStorage()

	assert.Equal(t, gaugeMap, memStrg.gaugeStorage)
	assert.Equal(t, counterMap, memStrg.counterStorage)
}

func Test_UpdateGaugeMetric(t *testing.T) {
	memStrg := MemStorage{
		gaugeStorage: make(map[string]float64),
	}

	mtrKey := "test-name"
	val := 12.5

	_ = memStrg.UpdateGaugeMetric(mtrKey, val)

	assert.Equal(t, val, memStrg.gaugeStorage[mtrKey])
}

func Test_UpdateCounterMetric(t *testing.T) {
	memStrg := MemStorage{
		counterStorage: make(map[string]int64),
	}

	// first call
	mtrKey := "test-name"
	val := 12

	_ = memStrg.UpdateCounterMetric(mtrKey, int64(val))

	assert.Equal(t, int64(val), memStrg.counterStorage[mtrKey])

	//sec call
	val2 := 2
	_ = memStrg.UpdateCounterMetric(mtrKey, int64(val2))

	assert.Equal(t, int64(val+val2), memStrg.counterStorage[mtrKey])
}
