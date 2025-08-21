package client

import models "github.com/sweetheart0330/metrics-alert/internal/model"

//go:generate mockgen -source=./interface.go -destination=./../mocks/mock_client.go
type IClient interface {
	SendGaugeMetric(m models.Metrics) error
	SendCounterMetric(m models.Metrics) error
}
