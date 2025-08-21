package handler

import (
	"github.com/sweetheart0330/metrics-alert/internal/service/contracts"
)

type Handler struct {
	metric contracts.MetricService
}

func NewHandler(metric contracts.MetricService) Handler {
	return Handler{metric: metric}
}
