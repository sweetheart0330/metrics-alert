package agent

import models "github.com/sweetheart0330/metrics-alert/internal/model"

//go:generate mockgen -source=./interface.go -destination=./../mocks/mock_agent.go
type MetricCollector interface {
	GetGauge() map[string]*float64
	GetCounter() models.Metrics
}
