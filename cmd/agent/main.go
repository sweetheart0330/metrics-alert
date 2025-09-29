package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/sweetheart0330/metrics-alert/internal/app"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			stop()
			return nil
		case <-egCtx.Done():
			return nil
		}
	})

	eg.Go(func() error { return app.RunAgent(egCtx) })

	if err := eg.Wait(); err != nil {
		log.Println("Error running agent", zap.Error(err))
		return
	}

	log.Println("Agent exited")
}

//
//func main() {
//	logger, err := zap.NewDevelopment()
//	if err != nil {
//		fmt.Println("failed to init logger, err: %w", err)
//		return
//	}
//	sugar := *logger.Sugar()
//
//	opt, err := getAgentFlags()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	clCfg := httpCl.Config{Host: "http://" + opt.Host}
//	cl := httpCl.NewClient(clCfg)
//	ag := agent.NewAgent(cl, nil, opt.ReportInterval, opt.PollInterval, &sugar)
//
//	for {
//		if err := ag.Run(); err != nil {
//			fmt.Println(err)
//		}
//		time.Sleep(1 * time.Second)
//	}
//}
//func getAgentFlags() (fl app.StartFlags, err error) {
//	err = env.Parse(&fl)
//	if err != nil {
//		return app.StartFlags{}, fmt.Errorf("failed to parse agent flags, err: %w", err)
//	}
//
//	if len(fl.Host) == 0 {
//		flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
//	}
//
//	if fl.ReportInterval == 0 {
//		flag.UintVar(&fl.ReportInterval, "r", 10, "interval between sending requests")
//	}
//
//	if fl.PollInterval == 0 {
//		flag.UintVar(&fl.PollInterval, "p", 2, "interval between collecting metrics")
//	}
//
//	flag.Parse()
//
//	return fl, nil
//}
