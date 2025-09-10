package app

import (
	"flag"
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
		return StartFlags{}, err
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

func getServerFlags() (fl string) {
	flag.StringVar(&fl, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return fl
}
