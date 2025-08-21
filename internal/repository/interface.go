package repository

type IRepository interface {
	UpdateCounterMetric(id string, val int64) error
	UpdateGaugeMetric(id string, val float64) error
}
