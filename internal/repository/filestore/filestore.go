package filestore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

type FileStorage struct {
	file *os.File // файл для записи
}

func NewFileStorage(filename string) (*FileStorage, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file, err: %w", err)
	}

	return &FileStorage{file: file}, nil
}

func (f *FileStorage) WriteMetrics(metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	jsonMetrics, err := json.MarshalIndent(metrics, "", "	")
	if err != nil {
		return err
	}

	if _, err = f.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err = f.file.Truncate(0); err != nil {
		return err
	}

	if _, err = f.file.Write(jsonMetrics); err != nil {
		return err
	}

	return f.file.Sync()
}

func (f *FileStorage) UploadMetrics() ([]models.Metrics, error) {
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f.file)
	if err != nil {
		return nil, err
	}

	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var metrics []models.Metrics
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}
