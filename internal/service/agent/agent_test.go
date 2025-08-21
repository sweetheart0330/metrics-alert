package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sweetheart0330/metrics-alert/internal/mocks"
	model "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/mock/gomock"
)

func Test_NewAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	mockCollector := mocks.NewMockMetricCollector(ctrl)

	ag := NewAgent(mockCl, mockCollector, 10)

	assert.Equal(t, mockCl, ag.cl)
	assert.Equal(t, mockCollector, ag.collect)
}
func Test_StartAgent(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	mockCollector := mocks.NewMockMetricCollector(ctrl)

	type args struct {
		ctx      context.Context
		cancel   context.CancelFunc
		gaugeMap map[string]*float64
		counter  model.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		prepare func(args args, err error)
	}{
		{
			name: "success",
			args: args{
				gaugeMap: map[string]*float64{
					"test1": new(float64),
					"test2": new(float64),
				},
				counter: model.Metrics{
					ID:    "test-counter",
					MType: model.Counter,
					Delta: ptrInt(5),
				},
			},
			wantErr: nil,
			prepare: func(args args, err error) {
				mockCollector.EXPECT().GetGauge().Return(args.gaugeMap)
				for k, v := range args.gaugeMap {
					mockCl.EXPECT().SendGaugeMetric(model.Metrics{
						ID:    k,
						MType: model.Gauge,
						Value: v,
					}).Return(nil)
				}

				mockCollector.EXPECT().GetCounter().Return(args.counter)
				mockCl.EXPECT().SendCounterMetric(args.counter).Return(nil)

				args.cancel() // завершаем работу цикла
			},
		},
		{
			name: "err in send gauge",
			args: args{
				gaugeMap: map[string]*float64{
					"test1": new(float64),
					"test2": new(float64),
				},
			},
			wantErr: errors.New("failed to send gauge request"),
			prepare: func(args args, err error) {
				mockCollector.EXPECT().GetGauge().Return(args.gaugeMap)
				mockCl.EXPECT().SendGaugeMetric(gomock.Any()).Return(err)
			},
		},
		{
			name: "failed to send counter",
			args: args{
				gaugeMap: map[string]*float64{
					"test1": new(float64),
					"test2": new(float64),
				},
				counter: model.Metrics{
					ID:    "test-counter",
					MType: model.Counter,
					Delta: ptrInt(5),
				},
			},
			wantErr: errors.New("failed to send counter request"),
			prepare: func(args args, err error) {
				mockCollector.EXPECT().GetGauge().Return(args.gaugeMap)
				for k, v := range args.gaugeMap {
					mockCl.EXPECT().SendGaugeMetric(model.Metrics{
						ID:    k,
						MType: model.Gauge,
						Value: v,
					}).Return(nil)
				}

				mockCollector.EXPECT().GetCounter().Return(args.counter)
				mockCl.EXPECT().SendCounterMetric(args.counter).Return(err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := Agent{cl: mockCl, collect: mockCollector, reportInterval: 10 * time.Second}
			tt.args.ctx, tt.args.cancel = context.WithCancel(context.Background())
			tt.prepare(tt.args, tt.wantErr)

			err := ag.StartAgent(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StartAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ptrInt(i int64) *int64 {
	return &i
}
