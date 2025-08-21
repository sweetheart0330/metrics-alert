package repository

//go:generate mockgen -source=./interface.go -destination=./../mocks/mock_repo.go
type IRepository interface {
	UpdateCounterMetric(id string, val int64) error
	UpdateGaugeMetric(id string, val float64) error
}
