package memory

import "fmt"

type MemStorage struct {
	gaugeStorage   map[string]float64
	counterStorage map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeStorage:   make(map[string]float64),
		counterStorage: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGaugeMetric(name string, value float64) error {
	ms.gaugeStorage[name] = value

	fmt.Println("current gauge value: ", value)
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(name string, value int64) error {
	val, ok := ms.counterStorage[name]
	if !ok {
		ms.counterStorage[name] = value
		fmt.Println("new counter value: ", ms.counterStorage[name])
		return nil
	}

	ms.counterStorage[name] = val + value

	fmt.Println("current counter value: ", ms.counterStorage[name])
	return nil
}
