package runtime

import (
	"context"
	"fmt"
	"runtime"
	"time"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

const (
	//gauge metrics
	AllocKey         = "Alloc"
	BuckHashSysKey   = "BuckHashSys"
	FreesKey         = "Frees"
	GCCPUFractionKey = "GCCPUFraction"
	GCSysKey         = "GCSys"
	HeapAllocKey     = "HeapAlloc"
	HeapIdleKey      = "HeapIdle"
	HeapInuseKey     = "HeapInuse"
	HeapObjectsKey   = "HeapObjects"
	HeapReleasedKey  = "HeapReleased"
	HeapSysKey       = "HeapSys"
	LastGCKey        = "LastGC"
	LookupsKey       = "Lookups"
	MCacheInuseKey   = "MCacheInuse"
	MCacheSysKey     = "MCacheSys"
	MSpanInuseKey    = "MSpanInuse"
	MSpanSysKey      = "MSpanSys"
	MallocsKey       = "Mallocs"
	NextGCKey        = "NextGC"
	NumForcedGCKey   = "NumForcedGC"
	NumGCKey         = "NumGC"
	OtherSysKey      = "OtherSys"
	PauseTotalNsKey  = "PauseTotalNs"
	StackInuseKey    = "StackInuse"
	StackSysKey      = "StackSys"
	SysKey           = "Sys"
	TotalAllocKey    = "TotalAlloc"

	//counter metrics
	PollCount = "PollCount"

	metricUpdateDelay = 2 * time.Second
)

type RuntimeMetrics struct {
	gauge   map[string]*float64
	counter int64
}

func NewRuntimeMetrics(ctx context.Context) *RuntimeMetrics {
	metric := &RuntimeMetrics{
		gauge: make(map[string]*float64),
	}

	go metric.startCollectMetrics(ctx)

	return metric
}

func (r *RuntimeMetrics) GetGauge() map[string]*float64 {
	return r.gauge
}

func (r *RuntimeMetrics) GetCounter() models.Metrics {
	return models.Metrics{
		ID:    PollCount,
		MType: models.Counter,
		Delta: &r.counter,
	}
}

func (r *RuntimeMetrics) startCollectMetrics(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r.collectMetrics()
			time.Sleep(metricUpdateDelay)
		}
	}
}

func (r *RuntimeMetrics) collectMetrics() map[string]*float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	floatPtr := func(val float64) *float64 {
		return &val
	}

	r.gauge[AllocKey] = floatPtr(float64(m.Alloc))
	r.gauge[BuckHashSysKey] = floatPtr(float64(m.BuckHashSys))
	r.gauge[FreesKey] = floatPtr(float64(m.Frees))
	r.gauge[GCCPUFractionKey] = floatPtr(m.GCCPUFraction)
	r.gauge[GCSysKey] = floatPtr(float64(m.GCSys))
	r.gauge[HeapAllocKey] = floatPtr(float64(m.HeapAlloc))
	r.gauge[HeapIdleKey] = floatPtr(float64(m.HeapIdle))
	r.gauge[HeapInuseKey] = floatPtr(float64(m.HeapInuse))
	r.gauge[HeapObjectsKey] = floatPtr(float64(m.HeapObjects))
	r.gauge[HeapReleasedKey] = floatPtr(float64(m.HeapReleased))
	r.gauge[HeapSysKey] = floatPtr(float64(m.HeapSys))
	r.gauge[LastGCKey] = floatPtr(float64(m.LastGC))
	r.gauge[LookupsKey] = floatPtr(float64(m.Lookups))
	r.gauge[MCacheInuseKey] = floatPtr(float64(m.MCacheInuse))
	r.gauge[MCacheSysKey] = floatPtr(float64(m.MCacheSys))
	r.gauge[MSpanInuseKey] = floatPtr(float64(m.MSpanInuse))
	r.gauge[MSpanSysKey] = floatPtr(float64(m.MSpanSys))
	r.gauge[MallocsKey] = floatPtr(float64(m.Mallocs))
	r.gauge[NextGCKey] = floatPtr(float64(m.NextGC))
	r.gauge[NumForcedGCKey] = floatPtr(float64(m.NumForcedGC))
	r.gauge[NumGCKey] = floatPtr(float64(m.NumGC))
	r.gauge[OtherSysKey] = floatPtr(float64(m.OtherSys))
	r.gauge[PauseTotalNsKey] = floatPtr(float64(m.PauseTotalNs))
	r.gauge[StackInuseKey] = floatPtr(float64(m.StackInuse))
	r.gauge[StackSysKey] = floatPtr(float64(m.StackSys))
	r.gauge[SysKey] = floatPtr(float64(m.Sys))
	r.gauge[TotalAllocKey] = floatPtr(float64(m.TotalAlloc))

	r.counter++

	fmt.Println("gauge collected")
	return r.gauge
}
