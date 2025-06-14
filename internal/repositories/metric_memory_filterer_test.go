package repositories

import (
	"context"
	"sync"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/require"
)

func int64Ptr(i int64) *int64 { return &i }

func TestMetricMemoryFiltererRepository_Filter(t *testing.T) {
	mu = &sync.RWMutex{}

	metric1 := types.Metrics{
		ID:    "metric1",
		MType: "gauge",
		Value: float64Ptr(100.1),
	}
	metric2 := types.Metrics{
		ID:    "metric2",
		MType: "counter",
		Delta: int64Ptr(10),
	}

	data = map[types.MetricID]types.Metrics{
		{ID: metric1.ID, MType: metric1.MType}: metric1,
		{ID: metric2.ID, MType: metric2.MType}: metric2,
	}

	repo := NewMetricMemoryFiltererRepository()

	// Test filtering existing metrics
	filterIDs := []types.MetricID{
		{ID: "metric1", MType: "gauge"},
		{ID: "metric2", MType: "counter"},
	}
	filtered, err := repo.Filter(context.Background(), filterIDs)
	require.NoError(t, err)
	require.Len(t, filtered, 2)
	require.Equal(t, metric1, filtered[filterIDs[0]])
	require.Equal(t, metric2, filtered[filterIDs[1]])

	// Test filtering with some non-existing IDs
	filterIDs = []types.MetricID{
		{ID: "metric1", MType: "gauge"},
		{ID: "metricX", MType: "gauge"},
	}
	filtered, err = repo.Filter(context.Background(), filterIDs)
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	require.Equal(t, metric1, filtered[filterIDs[0]])

	// Test filtering with empty slice
	filtered, err = repo.Filter(context.Background(), []types.MetricID{})
	require.NoError(t, err)
	require.Empty(t, filtered)
}
