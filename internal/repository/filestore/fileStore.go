package filestore

import (
	"encoding/json"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"io"
	"os"
)

type FileStorage struct {
	file *os.File // файл для записи
}

func NewFileStorage(filename string) (*FileStorage, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorage{file: file}, nil
}

func (f *FileStorage) WriteMetrics(metrics []models.Metrics) error {

	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	// Перемещаемся на начало
	if _, err = f.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	// Обрезаем файл до нулевой длины
	if err = f.file.Truncate(0); err != nil {
		return err
	}
	// Записываем новые данные
	if _, err = f.file.Write(jsonMetrics); err != nil {
		return err
	}
	// Опционально: чтобы гарантированно сбросить на диск
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
	// Возвращаем управление позицией на начало
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var metrics []models.Metrics
	err = json.Unmarshal(data, metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (f *FileStorage) Close() error {
	// закрываем файл
	return f.file.Close()
}
