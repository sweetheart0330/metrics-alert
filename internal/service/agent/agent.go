package agent

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/agent"
	"github.com/sweetheart0330/metrics-alert/internal/agent/runtime"
	"github.com/sweetheart0330/metrics-alert/internal/client"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/zap"
)

type Config struct {
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}
type Agent struct {
	cl      client.IClient
	collect agent.MetricCollector
	Config
	log            *zap.SugaredLogger
	counter        atomic.Int64
	reportInterval int64
	pollInterval   int64
}

func NewAgent(cl client.IClient, agent agent.MetricCollector, reportInterval uint, pollInterval uint, log *zap.SugaredLogger) *Agent {
	return &Agent{
		cl:      cl,
		collect: agent,
		Config: Config{
			ReportInterval: time.Duration(reportInterval) * time.Second,
			PollInterval:   time.Duration(pollInterval) * time.Second,
		},
		reportInterval: int64(reportInterval),
		pollInterval:   int64(pollInterval),
		log:            log,
	}
}

//func NewAgent(cl client.IClient, agent agent.MetricCollector, reportInterval, pollInterval uint, log *zap.SugaredLogger) *Agent {
//	return &Agent{
//		cl:      cl,
//		collect: agent,
//		Config:  Config{ReportInterval: time.Duration(reportInterval) * time.Second},
//		log:     log,
//
//		pollCount:      1,
//		counter:        1,
//		pollInterval:   int(pollInterval),
//		reportInterval: int(reportInterval),
//	}
//}

//	func (a *Agent) StartAgent(ctx context.Context) error {
//		for {
//			if err := a.Run(); err != nil {
//				fmt.Println(err)
//			}
//			time.Sleep(1 * time.Second)
//		}
//	}
//
//	func (a *Agent) Run() error {
//		var metrics []models.Metrics
//
//		if a.counter%a.pollInterval == 0 {
//			fmt.Println("collect metrics")
//			metrics = runtime.PullMetrics(a.pollCount)
//			a.pollCount++
//		}
//
//		if a.counter%a.reportInterval == 0 {
//			fmt.Println("report metrics")
//			for _, metric := range metrics {
//				if metric.MType == models.Gauge {
//					if err := a.cl.SendGaugeMetric(metric); err != nil {
//						return err
//					}
//				} else if metric.MType == models.Counter {
//					if err := a.cl.SendCounterMetric(metric); err != nil {
//						return err
//					}
//				}
//			}
//
//		}
//		a.counter++
//
//		return nil
//
// }
func (a *Agent) StartAgent(ctx context.Context) error {
	tick := time.NewTicker(a.ReportInterval)
	defer tick.Stop()
	tickP := time.NewTicker(a.PollInterval)
	defer tickP.Stop()

	var metrics []models.Metrics

	counter := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tickP.C:
			metrics = runtime.PullMetrics(int64(counter))
			counter++

			a.log.Info("metrics collected")
		case <-tick.C:
			err := a.sendNewMetrics(metrics)
			if err != nil {
				//a.log.Errorw("failed to send request to update metrics", "error", err.Error())
				//continue
				return fmt.Errorf("failed to send metrics: %w", err)
			}

			a.log.Info("metrics sent")
		}
	}
}

//func (a *Agent) StartAlterAgent(ctx context.Context) error {
//	//tick := time.NewTicker(a.ReportInterval)
//	//defer tick.Stop()
//	//tickP := time.NewTicker(a.PollInterval)
//	//defer tickP.Stop()
//
//	var metrics []models.Metrics
//
//	var pollCount int64 = 1
//	var counter int64 = 1
//	for {
//
//		select {
//		case <-ctx.Done():
//			return nil
//		default:
//			if counter%a.pollInterval == 0 {
//				metrics = runtime.PullMetrics(pollCount)
//				pollCount++
//				a.log.Info("metrics collected")
//			}
//
//			if counter%a.reportInterval == 0 {
//				fmt.Println("pol: ", pollCount)
//				if err := a.sendNewMetrics(metrics); err != nil {
//					return err
//				}
//			}
//			counter++
//			time.Sleep(1 * time.Second)
//		}
//	}
//}

//
//func (a *Agent) sendMetrics(models []models.Metrics) (err error) {
//	counter := a.collect.GetCounter()
//	err = a.cl.SendCounterMetric(counter)
//	if err != nil {
//		a.log.Errorw("failed to send counter",
//			"error", err.Error())
//		return fmt.Errorf("collect send counter metric failed: %w", err)
//	}
//
//	a.log.Debugw("sent counter metric",
//		"key", counter.ID,
//		"value", counter.Delta)
//	gaugeMap := a.collect.GetGauge()
//
//	gaugeMap.Range(func(key, value interface{}) bool {
//		valFl := value.(float64)
//		err = a.cl.SendGaugeMetric(models.Metrics{
//			ID:    key.(string),
//			MType: models.Gauge,
//			Value: &valFl,
//		})
//
//		if err != nil {
//			a.log.Errorw("failed to send gauge",
//				"key", key,
//				"value", valFl,
//				"error", err.Error())
//			return true
//		}
//
//		a.log.Debugw("sent gauge metric",
//			"key", key,
//			"value", valFl)
//
//		return true
//	})
//
//	if err != nil {
//		return fmt.Errorf("collect send gauge metric failed: %w", err)
//	}
//
//	return nil
//}

func (a *Agent) sendNewMetrics(metrics []models.Metrics) error {
	for _, m := range metrics {
		if m.MType == models.Gauge {
			err := a.cl.SendGaugeMetric(m)
			if err != nil {
				a.log.Error("failed to send gauge", zap.Error(err))
			}
		} else if m.MType == models.Counter {
			err := a.cl.SendCounterMetric(m)
			if err != nil {
				a.log.Warn("failed to send counter metric", zap.Error(err))
			}
		}
	}

	return nil
}
