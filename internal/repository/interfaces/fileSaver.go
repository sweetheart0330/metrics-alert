package interfaces

import models "github.com/sweetheart0330/metrics-alert/internal/model"

//go:generate mockgen -source=./fileSaver.go -destination=./../../mocks/mock_file_saver.go -package=mocks
type FileSaver interface {
	WriteMetrics(metrics []models.Metrics) error
	UploadMetrics() ([]models.Metrics, error)
}
