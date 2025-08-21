package app

import (
	"fmt"
	"net/http"

	"github.com/sweetheart0330/metrics-alert/internal/handler"
	"github.com/sweetheart0330/metrics-alert/internal/repository/memory"
	"github.com/sweetheart0330/metrics-alert/internal/router"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
)

func Run() error {
	inMemoryRepo := memory.NewMemStorage()
	MetricServ := metric.New(inMemoryRepo)
	h := handler.NewHandler(MetricServ)

	route := router.NewRouter(h)

	fmt.Println("Listening on :8080")
	return http.ListenAndServe(":8080", route)
}
