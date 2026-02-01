package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"time"

	rn "github.com/sweetheart0330/metrics-alert/internal/agent/runtime"
	"github.com/sweetheart0330/metrics-alert/internal/client/http/audit"
	httpCl "github.com/sweetheart0330/metrics-alert/internal/client/http/metric"
	"github.com/sweetheart0330/metrics-alert/internal/config"
	"github.com/sweetheart0330/metrics-alert/internal/observer"
	"github.com/sweetheart0330/metrics-alert/internal/repository/filestore"
	"github.com/sweetheart0330/metrics-alert/internal/repository/interfaces"
	"github.com/sweetheart0330/metrics-alert/internal/repository/postgre"
	"golang.org/x/sync/errgroup"

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
	clCfg := httpCl.Config{Host: "http://" + cfg.Host, SecretKey: cfg.SecretKey}
	cl := httpCl.NewClient(clCfg)
	ag := rn.NewRuntimeMetrics(ctx, cfg.PollInterval, &sugar)
	serv := servAgent.NewAgent(cl, ag, cfg.ReportInterval, cfg.PollInterval, &sugar)

	sugar.Info("Agent started")

	return serv.StartAgent(ctx)
}

func RunServer(ctx context.Context) error {
	cfg, err := config.GetServer()
	//cfg, err := config.NewServerConfig()
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

	publisher, err := addObserver(ctx, cfg, sugar)
	if err != nil {
		return fmt.Errorf("failed to init publisher, err: %w", err)
	}

	MetricServ, err := metric.New(repo, sugar, publisher)
	if err != nil {
		return fmt.Errorf("failed to init metric service, err: %w", err)
	}

	h, err := handler.NewHandler(MetricServ, sugar, cfg.SecretKey, cfg.RateLimit)
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

func addObserver(ctx context.Context, cfg config.ServerConfig, log zap.SugaredLogger) (observer.Publisher, error) {
	publisher := observer.NewAuditPublisher(&log)
	if len(cfg.AuditFile) != 0 {
		storage, err := filestore.NewAuditFileStorage(cfg.AuditFile)
		if err != nil {
			return nil, fmt.Errorf("failed to init audit file storage, err: %w", err)
		}

		publisher.Register(storage)
	}

	if len(cfg.AuditURL) != 0 {
		cl := audit.NewClient(cfg.AuditURL)
		publisher.Register(cl)
	}

	return publisher, nil
}

func chooseRepo(ctx context.Context, log *zap.SugaredLogger, cfg config.ServerConfig) (interfaces.IRepository, error) {
	if len(cfg.DBAddress) != 0 {
		db, err := postgre.NewDatabase(ctx, cfg.DBAddress, log)
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
