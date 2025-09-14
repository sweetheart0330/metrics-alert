package app

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type StartFlags struct {
	Host            string `env:"ADDRESS"`
	ReportInterval  uint   `env:"REPORT_INTERVAL"`
	PollInterval    uint   `env:"POLL_INTERVAL"`
	StoreInterval   *uint  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
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
		flag.UintVar(&fl.ReportInterval, "r", 10, "interval between sending requests")
	}

	if fl.PollInterval == 0 {
		flag.UintVar(&fl.PollInterval, "p", 2, "interval between collecting metrics")
	}

	flag.Parse()

	return fl, nil
}

func getServerFlags() (host StartFlags, err error) {
	fl := StartFlags{}
	err = env.Parse(&fl)
	if err != nil {
		return StartFlags{}, fmt.Errorf("failed to parse server flags, err: %w", err)
	}

	if len(fl.Host) == 0 {
		flag.StringVar(&fl.Host, "a", "localhost:8080", "address and port to send requests")
		flag.Parse()
	}

	if fl.StoreInterval == nil {
		var storeTime uint
		flag.UintVar(&storeTime, "i", 300, "frequency of storing metrics")
		fl.StoreInterval = &storeTime
	}

	if len(fl.FileStoragePath) == 0 {
		flag.StringVar(&fl.FileStoragePath, "f", "storage.txt", "file to save metrics")
	}

	if fl.Restore == false {
		flag.BoolVar(&fl.Restore, "r", false, "downloading metrics at the start from a file")
	}

	return fl, nil
}
