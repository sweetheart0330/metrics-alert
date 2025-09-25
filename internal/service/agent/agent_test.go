package agent

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sweetheart0330/metrics-alert/internal/mocks"
	model "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_NewAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	mockCollector := mocks.NewMockMetricCollector(ctrl)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("failed to init logger, err: %w", err)
		return
	}

	defer logger.Sync()
	sugar := *logger.Sugar()

	ag := NewAgent(mockCl, mockCollector, 10, &sugar)

	assert.Equal(t, mockCl, ag.cl)
	assert.Equal(t, mockCollector, ag.collect)
}
func Test_StartAgent(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	mockCollector := mocks.NewMockMetricCollector(ctrl)

	syncMap := sync.Map{}
	syncMap.Store("test1", 13.4)
	syncMap.Store("test2", 13.5)
	type args struct {
		ctx      context.Context
		cancel   context.CancelFunc
		gaugeMap *sync.Map
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
				gaugeMap: &syncMap,
				counter: model.Metrics{
					ID:    "test-counter",
					MType: model.Counter,
					Delta: ptrInt(5),
				},
			},
			wantErr: nil,
			prepare: func(args args, err error) {
				mockCollector.EXPECT().GetGauge().Return(args.gaugeMap)
				args.gaugeMap.Range(func(k, v interface{}) bool {
					fl := v.(float64)
					mockCl.EXPECT().SendGaugeMetric(model.Metrics{
						ID:    k.(string),
						MType: model.Gauge,
						Value: &fl,
					}).Return(nil)

					return true
				})

				mockCollector.EXPECT().GetCounter().Return(args.counter)
				mockCl.EXPECT().SendCounterMetric(args.counter).Return(nil)

				args.cancel() // завершаем работу цикла
			},
		},
		{
			name: "err in send gauge",
			args: args{
				gaugeMap: &syncMap,
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
				gaugeMap: &syncMap,
				counter: model.Metrics{
					ID:    "test-counter",
					MType: model.Counter,
					Delta: ptrInt(5),
				},
			},
			wantErr: errors.New("failed to send counter request"),
			prepare: func(args args, err error) {
				mockCollector.EXPECT().GetGauge().Return(args.gaugeMap)
				args.gaugeMap.Range(func(k, v interface{}) bool {
					fl := v.(float64)
					mockCl.EXPECT().SendGaugeMetric(model.Metrics{
						ID:    k.(string),
						MType: model.Gauge,
						Value: &fl,
					}).Return(nil)

					return true
				})

				mockCollector.EXPECT().GetCounter().Return(args.counter)
				mockCl.EXPECT().SendCounterMetric(args.counter).Return(err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := Agent{cl: mockCl, collect: mockCollector, Config: Config{ReportInterval: 1 * time.Second}}
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
