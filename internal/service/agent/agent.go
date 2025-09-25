package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/agent"
	"github.com/sweetheart0330/metrics-alert/internal/client"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/zap"
)

type Config struct {
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}
type Agent struct {
	cl      client.IClient
	collect agent.MetricCollector
	Config
	log *zap.SugaredLogger
}

func NewAgent(cl client.IClient, agent agent.MetricCollector, reportInterval uint, log *zap.SugaredLogger) *Agent {
	return &Agent{
		cl:      cl,
		collect: agent,
		Config:  Config{ReportInterval: time.Duration(reportInterval) * time.Second},
		log:     log,
	}
}

func (a Agent) StartAgent(ctx context.Context) error {
	tick := time.NewTicker(a.ReportInterval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context canceled")
			return nil
		case <-tick.C:
			err := a.sendMetrics()
			if err != nil {
				fmt.Println("err here:", err)
				return fmt.Errorf("failed to send metrics: %w", err)
			}

			a.log.Info("send metrics to server")
		}
	}
}

func (a Agent) sendMetrics() (err error) {
	gaugeMap := a.collect.GetGauge()

	gaugeMap.Range(func(key, value interface{}) bool {
		valFl := value.(float64)
		err = a.cl.SendGaugeMetric(models.Metrics{
			ID:    key.(string),
			MType: models.Gauge,
			Value: &valFl,
		})

		return err == nil
	})

	if err != nil {
		return fmt.Errorf("collect send gauge metric failed: %w", err)
	}

	counter := a.collect.GetCounter()
	err = a.cl.SendCounterMetric(counter)
	if err != nil {
		return fmt.Errorf("collect send counter metric failed: %w", err)
	}

	return nil
}
