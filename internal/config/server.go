package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Host            string `env:"ADDRESS"`
	StoreInterval   *uint  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DBAddress       string `env:"DATABASE_DSN"`
	SecretKey       string `env:"KEY"`
	RateLimit       int    `env:"RATE_LIMIT"`
	AuditFile       string `env:"AUDIT_FILE"`
	AuditURL        string `env:"AUDIT_URL"`
}

func GetServer() (host ServerConfig, err error) {
	fl := ServerConfig{}
	flag.StringVar(&fl.Host, "a", ":8080", "address and port to send requests")
	fl.StoreInterval = flag.Uint("i", 300, "frequency of storing metrics")
	flag.StringVar(&fl.FileStoragePath, "f", "storage.txt", "file to save metrics")
	flag.BoolVar(&fl.Restore, "r", false, "downloading metrics at the start from a file")
	flag.StringVar(&fl.DBAddress, "d", "", "downloading metrics at the start from a file")
	flag.StringVar(&fl.SecretKey, "k", "", "secret key")
	flag.IntVar(&fl.RateLimit, "l", 10, "rate limit")
	flag.StringVar(&fl.AuditFile, "observer-file", "", "observer file")
	flag.StringVar(&fl.AuditURL, "observer-url", "", "observer url")
	flag.Parse()

	err = env.Parse(&fl)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("failed to parse server flags, err: %w", err)
	}

	return fl, nil
}
