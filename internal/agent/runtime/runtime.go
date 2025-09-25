package runtime

import (
	"context"
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

func NewRuntimeMetrics(ctx context.Context, pollInterval uint, log *zap.SugaredLogger) *Metrics {
	metric := &Metrics{
		pollInterval: time.Duration(pollInterval) * time.Second,
		log:          log,
	}

	go metric.startCollectMetrics(ctx)

	return metric
}

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

	r.counter.Add(1)

	r.log.Info("gauge collected")
}
