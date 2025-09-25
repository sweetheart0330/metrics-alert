package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	servMetric "github.com/sweetheart0330/metrics-alert/internal/service/metric"
)

func (h Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := fillMetric(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.metric.UpdateMetric(*metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := h.getMetricFromBody(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get body, err: %v", err), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.metric.UpdateMetric(*metric)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update metric, err: %v", err), http.StatusBadRequest)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set(compressHeader, compressFormat)
	w.WriteHeader(http.StatusOK)
}

func (h Handler) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := h.getMetricFromBody(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get body, err: %v", err), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := h.metric.GetMetric(metric.ID)
	if err != nil {
		if errors.Is(err, servMetric.ErrMetricNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonMetric, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonMetric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(compressHeader, compressFormat)
	w.WriteHeader(http.StatusOK)
}

func (h Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metric := models.Metrics{}
	metric.MType = r.PathValue(models.TypeParam)
	if len(metric.MType) == 0 {
		http.Error(w, "failed to parse metric type from URL", http.StatusBadRequest)
		return
	}

	metric.ID = r.PathValue(models.NameParam)
	if len(metric.ID) == 0 {
		http.Error(w, "failed to parse metric name from URL", http.StatusBadRequest)
		return
	}

	resp, err := h.metric.GetMetric(metric.ID)
	if err != nil {
		if errors.Is(err, servMetric.ErrMetricNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch metric.MType {
	case models.Gauge:
		_, err = w.Write([]byte(strconv.FormatFloat(*resp.Value, 'f', -1, 64)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case models.Counter:
		_, err = w.Write([]byte(strconv.FormatInt(*resp.Delta, 10)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.metric.GetAllMetrics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Выполняем шаблон с данными метрик
	err = h.template.Execute(w, metrics)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона: "+err.Error(),
			http.StatusInternalServerError)
		return
	}
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

// Пользовательские функции для разыменования указателей
var funcMap = template.FuncMap{
	"derefInt": func(p *int64) int64 {
		if p != nil {
			return *p
		}
		return 0
	},
	"derefFloat": func(p *float64) float64 {
		if p != nil {
			return *p
		}
		return 0.0
	},
	"hasValue": func(p interface{}) bool {
		return p != nil
	},
}

func (h Handler) getMetricFromBody(w http.ResponseWriter, r *http.Request) (*models.Metrics, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, err: %w", err)
	}

	defer r.Body.Close()

	fmt.Println("body: ", string(body))
	var metric models.Metrics
	err = json.Unmarshal(body, &metric)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body, err: %w", err)
	}

	return &metric, nil
}
