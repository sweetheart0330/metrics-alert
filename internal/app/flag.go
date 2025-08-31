package app

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sweetheart0330/metrics-alert/internal/agent/runtime"
	myHTTP "github.com/sweetheart0330/metrics-alert/internal/client/http"
	"github.com/sweetheart0330/metrics-alert/internal/service/agent"
)

type Options struct {
	Client           myHTTP.Config
	Agent            agent.Config
	MetricsCollector runtime.Config
}

func getAgentOptions() (op Options, err error) {
	err = env.Parse(&op)
	if err != nil {
		return Options{}, err
	}

	var repInterval, pollInterval int
	if len(op.Client.Host) == 0 {
		flag.StringVar(&op.Client.Host, "a", "localhost:8080", "address and port to send requests")
	}

	if op.Agent.ReportInterval == 0 {
		flag.IntVar(&repInterval, "r", 10, "interval between sending requests")
		op.Agent.ReportInterval = time.Duration(repInterval) * time.Second
	}

	if op.MetricsCollector.PollInterval == 0 {
		flag.IntVar(&pollInterval, "p", 2, "interval between collecting metrics")
		op.MetricsCollector.PollInterval = time.Duration(pollInterval) * time.Second
	}

	flag.Parse()

	return op, nil
}

func getServerAddress() (fl string) {
	flag.StringVar(&fl, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return fl
}
