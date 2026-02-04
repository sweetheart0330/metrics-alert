package filestore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/observer"
)

type AuditFileStorage struct { // не стал объединять с другой структурой, решил, что у каждой структуры должна быть только одна задача
	file *os.File // файл для записи
}

func NewAuditFileStorage(path string) (observer.Observer, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file, err: %w", err)
	}

	return &AuditFileStorage{file: file}, nil
}

func (f *AuditFileStorage) Consume(ctx context.Context, ev models.AuditEvent) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	b = append(b, '\n')
	if _, err = f.file.Write(b); err != nil {
		return fmt.Errorf("failed to write audit to file, err: %w", err)
	}

	return f.file.Sync()
}

func (f *AuditFileStorage) Close() error {
	return f.file.Close()
}
