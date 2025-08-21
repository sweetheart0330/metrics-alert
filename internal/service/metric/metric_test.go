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

func Test_UpdateGaugeMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIRepository(ctrl)

	tests := []struct {
		name    string
		metric  models.Metrics
		wantErr error
	}{
		{
			name: "success",
			metric: models.Metrics{
				ID:    "test-metrics",
				Value: ptrFloat(5.2),
			},
			wantErr: nil,
		},
		{
			name: "error",
			metric: models.Metrics{
				Value: ptrFloat(0),
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{repo: mockRepo}

			mockRepo.EXPECT().UpdateGaugeMetric(tt.metric.ID, *tt.metric.Value).Return(tt.wantErr)

			err := m.UpdateGaugeMetric(tt.metric)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateGaugeMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_UpdateCounterMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIRepository(ctrl)

	tests := []struct {
		name    string
		metric  models.Metrics
		wantErr error
	}{
		{
			name: "success",
			metric: models.Metrics{
				ID:    "test-metrics",
				Delta: ptrInt(1),
			},
			wantErr: nil,
		},
		{
			name: "error",
			metric: models.Metrics{
				Delta: ptrInt(0),
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{repo: mockRepo}

			mockRepo.EXPECT().UpdateCounterMetric(tt.metric.ID, *tt.metric.Delta).Return(tt.wantErr)

			err := m.UpdateCounterMetric(tt.metric)
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
