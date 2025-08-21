package metric

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweetheart0330/metrics-alert/internal/mocks"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/mock/gomock"
)

func Test_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIRepository(ctrl)

	m := New(mockRepo)

	assert.Equal(t, mockRepo, m.repo)
}

func Test_UpdateMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIRepository(ctrl)

	tests := []struct {
		name    string
		metric  models.Metrics
		wantErr error
		prepare func(metric models.Metrics, err error)
	}{
		{
			name: "success counter",
			metric: models.Metrics{
				ID:    "test-metrics",
				MType: models.Counter,
				Delta: ptrInt(1),
			},
			wantErr: nil,
			prepare: func(metric models.Metrics, err error) {
				mockRepo.EXPECT().UpdateCounterMetric(metric).Return(err)
			},
		},
		{
			name: "undefined metric type",
			metric: models.Metrics{
				Delta: ptrInt(0),
			},
			wantErr: ErrUnknownMetricType,
			prepare: func(metric models.Metrics, err error) {},
		},
		{
			name: "success gauge",
			metric: models.Metrics{
				ID:    "test-metrics",
				MType: models.Gauge,
				Value: ptrFloat(5.2),
			},
			wantErr: nil,
			prepare: func(metric models.Metrics, err error) {
				mockRepo.EXPECT().UpdateGaugeMetric(metric).Return(err)
			},
		},
		{
			name: "error",
			metric: models.Metrics{
				MType: models.Gauge,
				Value: ptrFloat(0),
			},
			wantErr: errors.New("error"),
			prepare: func(metric models.Metrics, err error) {
				mockRepo.EXPECT().UpdateGaugeMetric(metric).Return(err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{repo: mockRepo}

			tt.prepare(tt.metric, tt.wantErr)

			err := m.UpdateMetric(tt.metric)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateGaugeMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ptrFloat(f float64) *float64 {
	return &f
}

func ptrInt(f int64) *int64 {
	return &f
}
