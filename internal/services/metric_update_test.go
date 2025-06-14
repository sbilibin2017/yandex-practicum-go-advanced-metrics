package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestAccumulateMetrics(t *testing.T) {
	delta1, delta2 := int64(1), int64(2)
	metrics := []types.Metrics{
		{ID: "metric1", MType: "counter", Delta: &delta1},
		{ID: "metric1", MType: "counter", Delta: &delta2},
		{ID: "metric2", MType: "gauge", Value: float64Ptr(3.14)},
	}

	result := accumulateMetrics(metrics)

	require.Len(t, result, 2)

	for _, m := range result {
		if m.MType == types.Counter {
			assert.Equal(t, int64(3), *m.Delta)
		}
		if m.MType == types.Gauge {
			assert.Equal(t, 3.14, *m.Value)
		}
	}
}

func TestUpdateMetrics(t *testing.T) {
	deltaExisting, deltaNew := int64(10), int64(5)
	existingMetric := types.Metrics{ID: "m1", MType: "counter", Delta: &deltaExisting}
	newMetric := types.Metrics{ID: "m1", MType: "counter", Delta: &deltaNew}

	metrics := []types.Metrics{newMetric}
	metricsMap := map[types.MetricID]types.Metrics{
		{ID: existingMetric.ID, MType: existingMetric.MType}: existingMetric,
	}

	updated := updateMetrics(metrics, metricsMap)
	require.Len(t, updated, 1)

	assert.Equal(t, int64(15), *updated[0].Delta)
}

func TestFilterMetricsByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilterer := NewMockMetricUpdateFilterer(ctrl)

	metrics := []types.Metrics{
		{ID: "m1", MType: "counter"},
		{ID: "m2", MType: "gauge"},
	}

	expectedMap := map[types.MetricID]types.Metrics{
		{ID: metrics[0].ID, MType: metrics[0].MType}: metrics[0],
		{ID: metrics[1].ID, MType: metrics[1].MType}: metrics[1],
	}

	mockFilterer.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return(expectedMap, nil).
		Times(1)

	result, err := filterMetricsByIDs(context.Background(), mockFilterer, metrics)
	require.NoError(t, err)
	assert.Equal(t, expectedMap, result)
}

func TestSaveMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricUpdateSaver(ctrl)

	metrics := []types.Metrics{
		{ID: "m1", MType: "counter"},
		{ID: "m2", MType: "gauge"},
	}

	gomock.InOrder(
		mockSaver.EXPECT().Send(gomock.Any(), metrics[0]).Return(nil),
		mockSaver.EXPECT().Send(gomock.Any(), metrics[1]).Return(nil),
		mockSaver.EXPECT().Save(gomock.Any()).Return(nil),
	)

	err := saveMetrics(context.Background(), mockSaver, metrics)
	require.NoError(t, err)
}

func TestSaveMetrics_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricUpdateSaver(ctrl)
	metrics := []types.Metrics{{ID: "m1", MType: "counter"}}

	mockSaver.EXPECT().Send(gomock.Any(), gomock.Any()).Return(errors.New("send error")).Times(1)

	err := saveMetrics(context.Background(), mockSaver, metrics)
	assert.Error(t, err)
}

func TestSaveMetrics_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricUpdateSaver(ctrl)
	metrics := []types.Metrics{{ID: "m1", MType: "counter"}}

	mockSaver.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockSaver.EXPECT().Save(gomock.Any()).Return(errors.New("save error")).Times(1)

	err := saveMetrics(context.Background(), mockSaver, metrics)
	assert.Error(t, err)
}

func TestMetricUpdateService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricUpdateSaver(ctrl)
	mockFilterer := NewMockMetricUpdateFilterer(ctrl)

	service := NewMetricUpdateService(mockSaver, mockFilterer)

	delta1 := int64(5)
	metrics := []types.Metrics{
		{ID: "m1", MType: "counter", Delta: &delta1},
	}

	filteredMap := map[types.MetricID]types.Metrics{
		{ID: metrics[0].ID, MType: metrics[0].MType}: metrics[0],
	}

	mockFilterer.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return(filteredMap, nil).
		Times(1)

	mockSaver.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(len(metrics))
	mockSaver.EXPECT().Save(gomock.Any()).Return(nil).Times(1)

	err := service.Update(context.Background(), metrics)
	require.NoError(t, err)
}

func TestMetricUpdateService_Update_FilterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricUpdateSaver(ctrl)
	mockFilterer := NewMockMetricUpdateFilterer(ctrl)

	service := NewMetricUpdateService(mockSaver, mockFilterer)

	mockFilterer.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("filter error")).
		Times(1)

	err := service.Update(context.Background(), []types.Metrics{})
	assert.Error(t, err)
}

func TestMetricUpdateService_Update_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := NewMockMetricUpdateSaver(ctrl)
	mockFilterer := NewMockMetricUpdateFilterer(ctrl)

	service := NewMetricUpdateService(mockSaver, mockFilterer)

	metrics := []types.Metrics{
		{ID: "m1", MType: "counter"},
	}

	mockFilterer.EXPECT().
		Filter(gomock.Any(), gomock.Any()).
		Return(map[types.MetricID]types.Metrics{}, nil).
		Times(1)

	mockSaver.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(len(metrics))
	mockSaver.EXPECT().Save(gomock.Any()).Return(errors.New("save error")).Times(1)

	err := service.Update(context.Background(), metrics)
	assert.Error(t, err)
}

func float64Ptr(f float64) *float64 {
	return &f
}
