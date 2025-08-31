package agent

import (
	"sync"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../mocks/mock_agent.go
type MetricCollector interface {
	GetGauge() *sync.Map
	GetCounter() models.Metrics
}
