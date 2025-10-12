package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Host            string `env:"ADDRESS"`
	StoreInterval   *uint  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DBAddress       string `env:"DATABASE_DSN"`
}

func NewServerConfig() (ServerConfig, error) {
	adrHost := flag.String("a", "localhost:8080", "Endpoint http server")
	storeInterval := flag.Int("i", 2, "interval save metrics in storage")
	fileStoragePath := flag.String("f", "storage.json", "filepath save metric storage")
	databaseDSN := flag.String("d", "", "database DSN")
	restore := flag.Bool("r", true, "read file storage metrics")
	flag.Parse()

	varAdrHost, ok := os.LookupEnv("ADDRESS")
	if ok {
		adrHost = &varAdrHost
	}
	varStoreInterval, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		intStoreInterval, err := strconv.Atoi(varStoreInterval)
		if err != nil {
			return ServerConfig{}, fmt.Errorf("error converting STORE_INTERVAL to int: %v", err)
		}
		storeInterval = &intStoreInterval
	}
	varFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		fileStoragePath = &varFileStoragePath
	}
	varRestore, ok := os.LookupEnv("RESTORE")
	if ok {
		boolVarRestore, err := strconv.ParseBool(varRestore)
		if err != nil {
			return ServerConfig{}, fmt.Errorf("error converting RESTORE to bool: %v", err)
		}
		restore = &boolVarRestore
	}
	varDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		databaseDSN = &varDatabaseDSN
	}

	someI := *storeInterval
	someUI := uint(someI)
	return ServerConfig{
		Host:            *adrHost,
		StoreInterval:   &someUI,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		DBAddress:       *databaseDSN,
	}, nil
}

func GetServer() (host ServerConfig, err error) {
	fl := ServerConfig{}
	flag.StringVar(&fl.Host, "a", ":8080", "address and port to send requests")
	fl.StoreInterval = flag.Uint("i", 300, "frequency of storing metrics")
	flag.StringVar(&fl.FileStoragePath, "f", "storage.txt", "file to save metrics")
	flag.BoolVar(&fl.Restore, "r", false, "downloading metrics at the start from a file")
	flag.StringVar(&fl.DBAddress, "d", "", "downloading metrics at the start from a file")
	flag.Parse()

	err = env.Parse(&fl)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("failed to parse server flags, err: %w", err)
	}

	return fl, nil
}
