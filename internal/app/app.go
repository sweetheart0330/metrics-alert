package app

import (
	"context"
	"fmt"
	"github.com/sweetheart0330/metrics-alert/internal/repository/filestore"
	"net/http"

	"github.com/sweetheart0330/metrics-alert/internal/agent/runtime"
	httpCl "github.com/sweetheart0330/metrics-alert/internal/client/http"
	"github.com/sweetheart0330/metrics-alert/internal/handler"
	"github.com/sweetheart0330/metrics-alert/internal/repository/memory"
	"github.com/sweetheart0330/metrics-alert/internal/router"
	servAgent "github.com/sweetheart0330/metrics-alert/internal/service/agent"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
	"go.uber.org/zap"
)

func RunAgent(ctx context.Context) error {
	opt, err := getAgentFlags()
	if err != nil {
		return err
	}

	clCfg := httpCl.Config{Host: "http://" + opt.Host}
	cl := httpCl.NewClient(clCfg)
	ag := runtime.NewRuntimeMetrics(ctx, opt.PollInterval)
	serv := servAgent.NewAgent(cl, ag, opt.ReportInterval)

	return serv.StartAgent(ctx)
}

func RunServer() error {
	ctx := context.Background()
	srvCfg, err := getServerFlags()
	if err != nil {
		return fmt.Errorf("failed to get server flags, err: %w", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to init logger, err: %w", err)
	}

	defer logger.Sync()
	sugar := *logger.Sugar()
	fileStorage, err := filestore.NewFileStorage(srvCfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("failed to init file storage, err: %w", err)
	}

	inMemoryRepo := memory.NewMemStorage()
	MetricServ := metric.New(ctx, inMemoryRepo, fileStorage, *srvCfg.StoreInterval, sugar)
	h, err := handler.NewHandler(MetricServ, sugar)
	if err != nil {
		return fmt.Errorf("failed to create new handler: %w", err)
	}

	route := router.NewRouter(h)

	sugar.Infow("Starting server", "srvCfg", srvCfg.Host)

	return http.ListenAndServe(srvCfg.Host, route)
}
