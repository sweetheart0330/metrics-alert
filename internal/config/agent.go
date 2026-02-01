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
	SecretKey      string `env:"KEY"`
}

func GetAgent() (AgentConfig, error) {
	fl := AgentConfig{}

	flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
	flag.UintVar(&fl.ReportInterval, "r", 10, "interval between sending requests")
	flag.UintVar(&fl.PollInterval, "p", 2, "interval between collecting metrics")
	flag.StringVar(&fl.SecretKey, "k", "", "secret key")
	flag.Parse()

	err := env.Parse(&fl)
	if err != nil {
		return AgentConfig{}, fmt.Errorf("failed to parse agent flags, err: %w", err)
	}

	return fl, nil
}
