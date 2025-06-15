package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestMetricUpdateService_Update_TableDriven(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		metrics []types.Metrics
	}

	ctx := context.Background()

	initialDelta := int64(10)
	existingDelta := int64(5)
	valueGauge := 42.0
	valueMem := 100.0
	wantErr := assert.AnError

	tests := []struct {
		name        string
		setupMocks  func(mockSaver *MockMetricUpdateSaver, mockGetter *MockMetricUpdateGetter)
		args        args
		expectedErr error
	}{
		{
			name: "Counter metric with existing metric updates delta and saves",
			setupMocks: func(mockSaver *MockMetricUpdateSaver, mockGetter *MockMetricUpdateGetter) {
				mockGetter.EXPECT().
					Get(ctx, types.MetricID{ID: "requests_total", MType: types.Counter}).
					Return(&types.Metrics{
						ID:    "requests_total",
						MType: types.Counter,
						Delta: &existingDelta,
					}, nil)
				expectedDelta := initialDelta + existingDelta
				mockSaver.EXPECT().
					Save(ctx, gomock.AssignableToTypeOf(types.Metrics{})).
					DoAndReturn(func(_ context.Context, m types.Metrics) error {
						assert.Equal(t, "requests_total", m.ID)
						assert.Equal(t, types.Counter, m.MType)
						assert.NotNil(t, m.Delta)
						assert.Equal(t, expectedDelta, *m.Delta)
						return nil
					})
			},
			args: args{metrics: []types.Metrics{{
				ID:    "requests_total",
				MType: types.Counter,
				Delta: &initialDelta,
			}}},
			expectedErr: nil,
		},
		{
			name: "Simple gauge metric saves successfully",
			setupMocks: func(mockSaver *MockMetricUpdateSaver, mockGetter *MockMetricUpdateGetter) {
				mockSaver.EXPECT().
					Save(ctx, types.Metrics{
						ID:    "cpu_usage",
						MType: types.Gauge,
						Value: &valueGauge,
					}).
					Return(nil)
			},
			args: args{metrics: []types.Metrics{{
				ID:    "cpu_usage",
				MType: types.Gauge,
				Value: &valueGauge,
			}}},
			expectedErr: nil,
		},
		{
			name: "Save returns error",
			setupMocks: func(mockSaver *MockMetricUpdateSaver, mockGetter *MockMetricUpdateGetter) {
				mockSaver.EXPECT().
					Save(ctx, types.Metrics{
						ID:    "memory_usage",
						MType: types.Gauge,
						Value: &valueMem,
					}).
					Return(wantErr)
			},
			args: args{metrics: []types.Metrics{{
				ID:    "memory_usage",
				MType: types.Gauge,
				Value: &valueMem,
			}}},
			expectedErr: wantErr,
		},
		{
			name: "Getter returns error on counter metric",
			setupMocks: func(mockSaver *MockMetricUpdateSaver, mockGetter *MockMetricUpdateGetter) {
				mockGetter.EXPECT().
					Get(ctx, types.MetricID{ID: "requests_total", MType: types.Counter}).
					Return(nil, wantErr)
				// Save should not be called
			},
			args: args{metrics: []types.Metrics{{
				ID:    "requests_total",
				MType: types.Counter,
				Delta: &initialDelta,
			}}},
			expectedErr: wantErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver := NewMockMetricUpdateSaver(ctrl)
			mockGetter := NewMockMetricUpdateGetter(ctrl)

			service := NewMetricUpdateService(mockSaver, mockGetter)

			tt.setupMocks(mockSaver, mockGetter)

			err := service.Update(ctx, tt.args.metrics)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
