package runtime

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/zap"
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
	RandomValue      = "RandomValue"
	//counter metrics
	PollCount = "PollCount"
)

type Config struct {
	PollInterval time.Duration
}
type Metrics struct {
	gauge        sync.Map
	mu           sync.RWMutex
	counter      atomic.Int64
	pollInterval time.Duration
	log          *zap.SugaredLogger
}

//func NewRuntimeMetrics(ctx context.Context, pollInterval uint, log *zap.SugaredLogger) *Metrics {
//	metric := &Metrics{
//		pollInterval: time.Duration(pollInterval) * time.Second,
//		log:          log,
//	}
//
//	go metric.startCollectMetrics(ctx)
//
//	return metric
//}

func (r *Metrics) GetGauge() *sync.Map {
	r.mu.Lock()
	defer r.mu.Unlock()

	return &r.gauge
}

func (r *Metrics) GetCounter() models.Metrics {
	counter := r.counter.Load()
	return models.Metrics{
		ID:    PollCount,
		MType: models.Counter,
		Delta: &counter,
	}
}

func (r *Metrics) startCollectMetrics(ctx context.Context) {
	t := time.NewTicker(r.pollInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			r.collectMetrics()
		}
	}
}

func (r *Metrics) collectMetrics() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	r.gauge.Store(AllocKey, float64(m.Alloc))
	r.gauge.Store(BuckHashSysKey, float64(m.BuckHashSys))
	r.gauge.Store(FreesKey, float64(m.Frees))
	r.gauge.Store(GCCPUFractionKey, m.GCCPUFraction)
	r.gauge.Store(GCSysKey, float64(m.GCSys))
	r.gauge.Store(HeapAllocKey, float64(m.HeapAlloc))
	r.gauge.Store(HeapIdleKey, float64(m.HeapIdle))
	r.gauge.Store(HeapInuseKey, float64(m.HeapInuse))
	r.gauge.Store(HeapObjectsKey, float64(m.HeapObjects))
	r.gauge.Store(HeapReleasedKey, float64(m.HeapReleased))
	r.gauge.Store(HeapSysKey, float64(m.HeapSys))
	r.gauge.Store(LastGCKey, float64(m.LastGC))
	r.gauge.Store(LookupsKey, float64(m.Lookups))
	r.gauge.Store(MCacheInuseKey, float64(m.MCacheInuse))
	r.gauge.Store(MCacheSysKey, float64(m.MCacheSys))
	r.gauge.Store(MSpanInuseKey, float64(m.MSpanInuse))
	r.gauge.Store(MSpanSysKey, float64(m.MSpanSys))
	r.gauge.Store(MallocsKey, float64(m.Mallocs))
	r.gauge.Store(NextGCKey, float64(m.NextGC))
	r.gauge.Store(NumForcedGCKey, float64(m.NumForcedGC))
	r.gauge.Store(NumGCKey, float64(m.NumGC))
	r.gauge.Store(OtherSysKey, float64(m.OtherSys))
	r.gauge.Store(PauseTotalNsKey, float64(m.PauseTotalNs))
	r.gauge.Store(StackInuseKey, float64(m.StackInuse))
	r.gauge.Store(StackSysKey, float64(m.StackSys))
	r.gauge.Store(SysKey, float64(m.Sys))
	r.gauge.Store(TotalAllocKey, float64(m.TotalAlloc))
	r.gauge.Store(RandomValue, rand.Float64())

	r.counter.Add(1)

	r.log.Info("gauge collected")
}

func PullMetrics(pollCount int64) []models.Metrics {
	metRuntime := runtime.MemStats{}

	runtime.ReadMemStats(&metRuntime)

	metrics := []models.Metrics{
		{ID: "Alloc", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Alloc))},
		{ID: "BuckHashSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.BuckHashSys))},
		{ID: "Frees", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Frees))},
		{ID: "GCCPUFraction", MType: models.Gauge, Value: float64Ptr(metRuntime.GCCPUFraction)},
		{ID: "GCSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.GCSys))},
		{ID: "HeapAlloc", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapAlloc))},
		{ID: "HeapIdle", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapIdle))},
		{ID: "HeapInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapInuse))},
		{ID: "HeapObjects", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapObjects))},
		{ID: "HeapReleased", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapReleased))},
		{ID: "HeapSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.HeapSys))},
		{ID: "LastGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.LastGC))},
		{ID: "Lookups", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Lookups))},
		{ID: "MCacheInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MCacheInuse))},
		{ID: "MCacheSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MCacheSys))},
		{ID: "MSpanInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MSpanInuse))},
		{ID: "MSpanSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.MSpanSys))},
		{ID: "Mallocs", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Mallocs))},
		{ID: "NextGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.NextGC))},
		{ID: "NumForcedGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.NumForcedGC))},
		{ID: "NumGC", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.NumGC))},
		{ID: "OtherSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.OtherSys))},
		{ID: "PauseTotalNs", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.PauseTotalNs))},
		{ID: "StackInuse", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.StackInuse))},
		{ID: "StackSys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.StackSys))},
		{ID: "Sys", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.Sys))},
		{ID: "TotalAlloc", MType: models.Gauge, Value: float64Ptr(float64(metRuntime.TotalAlloc))},
		{ID: "RandomValue", MType: models.Gauge, Value: float64Ptr(rand.Float64())},
		{ID: "PollCount", MType: models.Counter, Delta: &pollCount},
	}
	return metrics
}

func float64Ptr(f float64) *float64 {
	return &f
}
