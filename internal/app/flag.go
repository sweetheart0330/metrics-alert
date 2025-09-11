package app

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type StartFlags struct {
	Host           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func getAgentFlags() (fl StartFlags, err error) {
	err = env.Parse(&fl)
	if err != nil {
		return StartFlags{}, fmt.Errorf("failed to parse agent flags, err: %w", err)
	}

	if len(fl.Host) == 0 {
		flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
	}

	if fl.ReportInterval == 0 {
		flag.IntVar(&fl.ReportInterval, "r", 10, "interval between sending requests")
	}

	if fl.PollInterval == 0 {
		flag.IntVar(&fl.PollInterval, "p", 2, "interval between collecting metrics")
	}

	flag.Parse()

	return fl, nil
}

func getServerFlags() (host string, err error) {
	fl := StartFlags{}
	err = env.Parse(&fl)
	if err != nil {
		return "", fmt.Errorf("failed to parse server flags, err: %w", err)
	}

	if len(fl.Host) == 0 {
		flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
		flag.Parse()
	}

	return fl.Host, nil
}
