package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/repository/filestore"
	"golang.org/x/sync/errgroup"

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
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to init logger, err: %w", err)
	}

	defer logger.Sync()
	sugar := *logger.Sugar()
	clCfg := httpCl.Config{Host: "http://" + opt.Host}
	cl := httpCl.NewClient(clCfg)
	ag := runtime.NewRuntimeMetrics(ctx, opt.PollInterval, &sugar)
	serv := servAgent.NewAgent(cl, ag, opt.ReportInterval, &sugar)

	return serv.StartAgent(ctx)
}

func RunServer(ctx context.Context) error {
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

	eg, egCtx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr:    srvCfg.Host,
		Handler: route,
	}
	eg.Go(func() error {
		sugar.Infow("Starting server", "srvCfg", srvCfg.Host)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sugar.Infow("Stopping server", "srvCfg", srvCfg.Host)

		return server.Shutdown(shCtx)
	})

	return eg.Wait()
}
