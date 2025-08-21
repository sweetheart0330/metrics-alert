package handler

import (
	"fmt"
	"net/http"
	"strconv"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

func (h Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := fillMetric(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case models.Counter:
		err = h.metric.UpdateCounterMetric(*metric)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case models.Gauge:
		err = h.metric.UpdateGaugeMetric(*metric)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func fillMetric(r *http.Request) (*models.Metrics, error) {
	metric := models.Metrics{}
	metric.MType = r.PathValue(models.TypeParam)
	if len(metric.MType) == 0 {
		return nil, fmt.Errorf("can't parse metric type from path")
	}

	metric.ID = r.PathValue(models.NameParam)
	if len(metric.ID) == 0 {
		return nil, fmt.Errorf("can't parse metric name")
	}

	rawVal := r.PathValue(models.ValueParam)
	if len(rawVal) == 0 {
		return nil, fmt.Errorf("can't parse metric value")
	}

	switch metric.MType {
	case models.Counter:
		delta, err := strconv.ParseInt(rawVal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse metric counter delta, err %w", err)
		}

		metric.Delta = &delta

	case models.Gauge:
		val, err := strconv.ParseFloat(rawVal, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse metric gauge value, err: %w", err)
		}

		metric.Value = &val
	default:
		return nil, fmt.Errorf("undefined metric type")
	}

	return &metric, nil
}
