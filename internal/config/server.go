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
}

func GetServer() (host ServerConfig, err error) {
	fl := ServerConfig{}
	err = env.Parse(&fl)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("failed to parse server flags, err: %w", err)
	}

	if len(fl.Host) == 0 {
		flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
	}

	if fl.StoreInterval == nil {
		fl.StoreInterval = flag.Uint("i", 300, "frequency of storing metrics")
	}

	if len(fl.FileStoragePath) == 0 {
		flag.StringVar(&fl.FileStoragePath, "f", "storage.txt", "file to save metrics")
	}

	if !fl.Restore {
		flag.BoolVar(&fl.Restore, "r", false, "downloading metrics at the start from a file")
	}

	flag.Parse()

	return fl, nil
}
