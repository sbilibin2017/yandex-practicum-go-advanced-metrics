package repositories

import (
	"context"
	"sync"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/require"
)

func float64Ptr(f float64) *float64 { return &f }

func TestMetricMemorySaverRepository_SendStoresMetric(t *testing.T) {
	mu = &sync.RWMutex{}
	data = make(map[types.MetricID]types.Metrics)
	repo := NewMetricMemorySaverRepository()

	metric := types.Metrics{
		ID:    "testMetric",
		MType: "gauge",
		Value: float64Ptr(123.45),
		Delta: nil,
	}

	err := repo.Send(context.Background(), metric)
	require.NoError(t, err)

	mu.RLock()
	defer mu.RUnlock()
	key := types.MetricID{ID: metric.ID, MType: metric.MType}
	stored, ok := data[key]
	require.True(t, ok)
	require.Equal(t, metric.ID, stored.ID)
	require.Equal(t, metric.MType, stored.MType)
	require.NotNil(t, stored.Value)
	require.Equal(t, *metric.Value, *stored.Value)
	require.Nil(t, stored.Delta)
}

func TestMetricMemorySaverRepository_SendStoresCounterMetric(t *testing.T) {
	mu := &sync.RWMutex{}
	data := make(map[types.MetricID]types.Metrics)
	repo := &MetricMemorySaverRepository{
		mu:   mu,
		data: data,
	}

	delta := int64(42)
	metric := types.Metrics{
		ID:    "counterMetric",
		MType: "counter",
		Delta: &delta,
		Value: nil,
	}

	err := repo.Send(context.Background(), metric)
	require.NoError(t, err)

	mu.RLock()
	defer mu.RUnlock()
	key := types.MetricID{ID: metric.ID, MType: metric.MType}
	stored, ok := data[key]
	require.True(t, ok)
	require.Equal(t, metric.ID, stored.ID)
	require.Equal(t, metric.MType, stored.MType)
	require.NotNil(t, stored.Delta)
	require.Equal(t, *metric.Delta, *stored.Delta)
	require.Nil(t, stored.Value)
}

func TestMetricMemorySaverRepository_SaveNoOp(t *testing.T) {
	mu := &sync.RWMutex{}
	data := make(map[types.MetricID]types.Metrics)
	repo := &MetricMemorySaverRepository{
		mu:   mu,
		data: data,
	}

	err := repo.Save(context.Background())
	require.NoError(t, err)
}
