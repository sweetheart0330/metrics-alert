package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sweetheart0330/metrics-alert/internal/agent/runtime"
	httpCl "github.com/sweetheart0330/metrics-alert/internal/client/http"
	"github.com/sweetheart0330/metrics-alert/internal/handler"
	"github.com/sweetheart0330/metrics-alert/internal/repository/memory"
	"github.com/sweetheart0330/metrics-alert/internal/router"
	servAgent "github.com/sweetheart0330/metrics-alert/internal/service/agent"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
)

const (
	host = "localhost:8080"
)

func RunAgent() error {
	ctx := context.Background()
	clCfg := httpCl.Config{Host: "http://" + host}
	cl := httpCl.NewClient(clCfg)
	ag := runtime.NewRuntimeMetrics(ctx)

	serv := servAgent.NewAgent(cl, ag)

	return serv.StartAgent(ctx)
}

func RunServer() error {
	inMemoryRepo := memory.NewMemStorage()
	MetricServ := metric.New(inMemoryRepo)
	h := handler.NewHandler(MetricServ)

	route := router.NewRouter(h)

	fmt.Println("Listening on ", host)
	return http.ListenAndServe(":8080", route)
}
