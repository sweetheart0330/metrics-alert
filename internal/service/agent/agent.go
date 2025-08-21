package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/agent"
	"github.com/sweetheart0330/metrics-alert/internal/client"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

type Agent struct {
	cl             client.IClient
	collect        agent.MetricCollector
	reportInterval time.Duration
}

func NewAgent(cl client.IClient, agent agent.MetricCollector, reportInterval time.Duration) *Agent {
	return &Agent{
		cl:             cl,
		collect:        agent,
		reportInterval: reportInterval,
	}
}

func (a Agent) StartAgent(ctx context.Context) error {
	tick := time.NewTicker(a.reportInterval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			err := a.sendMetrics()
			if err != nil {
				return fmt.Errorf("failed to send metrics: %w", err)
			}

			fmt.Println("metrics sent to server")
		}
	}
}

func (a Agent) sendMetrics() error {
	gaugeMap := a.collect.GetGauge()
	for k, v := range gaugeMap {
		err := a.cl.SendGaugeMetric(models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: v,
		})
		if err != nil {
			return fmt.Errorf("collect send gauge metric failed: %w", err)
		}
	}

	counter := a.collect.GetCounter()
	err := a.cl.SendCounterMetric(counter)
	if err != nil {
		return fmt.Errorf("collect send counter metric failed: %w", err)
	}

	return nil
}
