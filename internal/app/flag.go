package app

import (
	"flag"
	"time"

	"github.com/sweetheart0330/metrics-alert/internal/agent/runtime"
	myHTTP "github.com/sweetheart0330/metrics-alert/internal/client/http"
	"github.com/sweetheart0330/metrics-alert/internal/service/agent"
)

type Options struct {
	Client           myHTTP.Config
	Agent            agent.Config
	MetricsCollector runtime.Config
}

func getAgentOptions() (op Options) {
	var repInterval, pollInterval int
	flag.StringVar(&op.Client.Host, "a", "localhost:8080", "address and port to send requests")
	flag.IntVar(&repInterval, "r", 10, "interval between sending requests")
	flag.IntVar(&pollInterval, "p", 2, "interval between collecting metrics")
	flag.Parse()

	op.Agent.ReportInterval = time.Duration(repInterval) * time.Second
	op.MetricsCollector.PollInterval = time.Duration(pollInterval) * time.Second

	return op
}

func getServerAddress() (fl string) {
	flag.StringVar(&fl, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return fl
}
