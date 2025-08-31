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

func RunAgent(ctx context.Context) error {
	opt := getAgentOptions()

	clCfg := httpCl.Config{Host: "http://" + opt.Client.Host}
	cl := httpCl.NewClient(clCfg)
	ag := runtime.NewRuntimeMetrics(ctx, opt.MetricsCollector.PollInterval)
	serv := servAgent.NewAgent(cl, ag, opt.Agent.ReportInterval)

	return serv.StartAgent(ctx)
}

func RunServer() error {

	addr := getServerAddress()
	inMemoryRepo := memory.NewMemStorage()
	MetricServ := metric.New(inMemoryRepo)
	h, err := handler.NewHandler(MetricServ)
	if err != nil {
		return fmt.Errorf("failed to create new handler: %w", err)
	}

	route := router.NewRouter(h)

	fmt.Println("Listening on ", addr)
	return http.ListenAndServe(addr, route)
}
