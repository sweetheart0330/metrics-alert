package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/config"
	"github.com/sweetheart0330/metrics-alert/internal/repository/filestore"
	"github.com/sweetheart0330/metrics-alert/internal/repository/interfaces"
	"github.com/sweetheart0330/metrics-alert/internal/repository/postgre"
	"golang.org/x/sync/errgroup"

	httpCl "github.com/sweetheart0330/metrics-alert/internal/client/http"
	"github.com/sweetheart0330/metrics-alert/internal/handler"
	"github.com/sweetheart0330/metrics-alert/internal/repository/memory"
	"github.com/sweetheart0330/metrics-alert/internal/router"
	servAgent "github.com/sweetheart0330/metrics-alert/internal/service/agent"
	"github.com/sweetheart0330/metrics-alert/internal/service/metric"
	"go.uber.org/zap"
)

func RunAgent(ctx context.Context) error {
	cfg, err := config.GetAgent()
	if err != nil {
		return err
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to init logger, err: %w", err)
	}

	defer logger.Sync()
	sugar := *logger.Sugar()
	clCfg := httpCl.Config{Host: "http://" + cfg.Host}
	cl := httpCl.NewClient(clCfg)
	//ag := runtime.NewRuntimeMetrics(ctx, cfg.PollInterval, &sugar)
	serv := servAgent.NewAgent(cl, nil, cfg.ReportInterval, cfg.PollInterval, &sugar)

	sugar.Info("Agent started")

	return serv.StartAgent(ctx)
}

func RunServer(ctx context.Context) error {
	cfg, err := config.GetServer()
	if err != nil {
		return fmt.Errorf("failed to get server flags, err: %w", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to init logger, err: %w", err)
	}

	defer logger.Sync()
	sugar := *logger.Sugar()

	repo, err := chooseRepo(ctx, &sugar, cfg)
	if err != nil {
		return fmt.Errorf("failed to init repo, err: %w", err)
	}

	MetricServ, err := metric.New(repo, sugar)
	if err != nil {
		return fmt.Errorf("failed to init metric service, err: %w", err)
	}

	h, err := handler.NewHandler(MetricServ, sugar)
	if err != nil {
		return fmt.Errorf("failed to create new handler: %w", err)
	}

	route := router.NewRouter(h)

	eg, egCtx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: route,
	}
	eg.Go(func() error {
		sugar.Infow("Starting server", "cfg", cfg.Host)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sugar.Infow("Stopping server", "cfg", cfg.Host)

		return server.Shutdown(shCtx)
	})

	return eg.Wait()
}

func chooseRepo(ctx context.Context, log *zap.SugaredLogger, cfg config.ServerConfig) (interfaces.IRepository, error) {
	if len(cfg.DBAddress) != 0 {
		db, err := postgre.NewDatabase(ctx, cfg.DBAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to init database, err: %w", err)
		}

		return db, nil
	}

	fileStorage, err := filestore.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init file storage, err: %w", err)
	}

	inMemoryRepo := memory.NewMemStorage(ctx, fileStorage, log, cfg.Restore, *cfg.StoreInterval)

	return inMemoryRepo, nil
}
