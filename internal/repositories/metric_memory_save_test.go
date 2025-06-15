package repositories

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemorySaveRepository_Save(t *testing.T) {
	storage := engines.NewMemoryStorage[types.MetricID, types.Metrics]()
	repo := NewMetricMemorySaveRepository(storage)

	metric := types.Metrics{
		ID:    "metric1",
		MType: "gauge",
	}

	err := repo.Save(context.Background(), metric)
	assert.NoError(t, err)

	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	key := types.MetricID{ID: metric.ID, MType: metric.MType}

	savedMetric, ok := storage.Data[key]
	assert.True(t, ok, "metric should be saved")
	assert.Equal(t, metric.ID, savedMetric.ID)
	assert.Equal(t, metric.MType, savedMetric.MType)
}

func TestMetricMemorySaveRepository_ConcurrentSave(t *testing.T) {
	storage := engines.NewMemoryStorage[types.MetricID, types.Metrics]()
	repo := NewMetricMemorySaveRepository(storage)

	numMetrics := 100
	var wg sync.WaitGroup
	wg.Add(numMetrics)

	for i := 0; i < numMetrics; i++ {
		go func(i int) {
			defer wg.Done()
			metric := types.Metrics{
				ID:    fmt.Sprintf("metric%d", i),
				MType: "gauge",
			}
			err := repo.Save(context.Background(), metric)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	assert.Len(t, storage.Data, numMetrics)
}
