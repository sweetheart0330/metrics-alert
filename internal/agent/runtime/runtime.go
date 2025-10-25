package runtime

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	util "github.com/shirou/gopsutil/v4/mem"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/zap"
)

const (
	//metricMap metrics
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
	TotalMemory      = "TotalMemory"
	FreeMemory       = "FreeMemory"
	CPUutilization1  = "CPUutilization1"
	//counter metrics
	PollCount = "PollCount"
)

type Config struct {
	PollInterval time.Duration
}
type Metrics struct {
	metricMap    sync.Map
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

func (r *Metrics) GetMetrics() *sync.Map {
	return &r.metricMap
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
			r.collectUtilMetrics()
			r.collectMetrics()

			r.log.Debug("metrics collected")
		}
	}
}

func (r *Metrics) collectUtilMetrics() {
	v, _ := util.VirtualMemory()

	r.storeGaugeMetric(TotalMemory, float64(v.Total))
	r.storeGaugeMetric(FreeMemory, float64(v.Free))

	percentages, _ := cpu.Percent(time.Second, false)
	r.storeGaugeMetric(CPUutilization1, percentages[0])
}

func (r *Metrics) collectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	r.storeGaugeMetric(AllocKey, float64(m.Alloc))
	r.storeGaugeMetric(BuckHashSysKey, float64(m.BuckHashSys))
	r.storeGaugeMetric(FreesKey, float64(m.Frees))
	r.storeGaugeMetric(GCCPUFractionKey, float64(m.GCCPUFraction))
	r.storeGaugeMetric(GCSysKey, float64(m.GCSys))
	r.storeGaugeMetric(HeapAllocKey, float64(m.HeapAlloc))
	r.storeGaugeMetric(HeapIdleKey, float64(m.HeapIdle))
	r.storeGaugeMetric(HeapInuseKey, float64(m.HeapInuse))
	r.storeGaugeMetric(HeapObjectsKey, float64(m.HeapObjects))
	r.storeGaugeMetric(HeapReleasedKey, float64(m.HeapReleased))
	r.storeGaugeMetric(HeapSysKey, float64(m.HeapSys))
	r.storeGaugeMetric(LastGCKey, float64(m.LastGC))
	r.storeGaugeMetric(LookupsKey, float64(m.Lookups))
	r.storeGaugeMetric(MCacheInuseKey, float64(m.MCacheInuse))
	r.storeGaugeMetric(MCacheSysKey, float64(m.MCacheSys))
	r.storeGaugeMetric(MSpanInuseKey, float64(m.MSpanInuse))
	r.storeGaugeMetric(MSpanSysKey, float64(m.MSpanSys))
	r.storeGaugeMetric(MallocsKey, float64(m.Mallocs))
	r.storeGaugeMetric(NextGCKey, float64(m.NextGC))
	r.storeGaugeMetric(NumForcedGCKey, float64(m.NumForcedGC))
	r.storeGaugeMetric(NumGCKey, float64(m.NumGC))
	r.storeGaugeMetric(OtherSysKey, float64(m.OtherSys))
	r.storeGaugeMetric(PauseTotalNsKey, float64(m.PauseTotalNs))
	r.storeGaugeMetric(StackInuseKey, float64(m.StackInuse))
	r.storeGaugeMetric(StackSysKey, float64(m.StackSys))
	r.storeGaugeMetric(SysKey, float64(m.Sys))
	r.storeGaugeMetric(TotalAllocKey, float64(m.TotalAlloc))
	r.storeGaugeMetric(RandomValue, rand.Float64())

	newCount := r.counter.Add(1)
	r.metricMap.Store(PollCount, models.Metrics{
		ID:    PollCount,
		MType: models.Counter,
		Delta: &newCount,
	})
}

func (r *Metrics) storeGaugeMetric(key string, value float64) {
	r.metricMap.Store(key, models.Metrics{
		ID:    key,
		MType: models.Gauge,
		Value: &value,
	})
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
