package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Host           string `env:"ADDRESS"`
	ReportInterval uint   `env:"REPORT_INTERVAL"`
	PollInterval   uint   `env:"POLL_INTERVAL"`
}

func GetAgent() (fl AgentConfig, err error) {
	err = env.Parse(&fl)
	if err != nil {
		return AgentConfig{}, fmt.Errorf("failed to parse agent flags, err: %w", err)
	}

	if len(fl.Host) == 0 {
		flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
	}

	if fl.ReportInterval == 0 {
		flag.UintVar(&fl.ReportInterval, "r", 10, "interval between sending requests")
	}

	if fl.PollInterval == 0 {
		flag.UintVar(&fl.PollInterval, "p", 2, "interval between collecting metrics")
	}

	flag.Parse()

	return fl, nil
}
