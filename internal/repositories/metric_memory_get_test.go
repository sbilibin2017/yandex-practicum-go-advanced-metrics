package repositories

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryGetRepository_Get(t *testing.T) {
	storage := engines.NewMemoryStorage[types.MetricID, types.Metrics]()

	// Pre-fill storage with a metric
	existingMetric := types.Metrics{
		ID:    "metric1",
		MType: "gauge",
		// You can set Value, Delta, etc. if required
	}

	key := types.MetricID{ID: existingMetric.ID, MType: existingMetric.MType}

	storage.Mu.Lock()
	storage.Data[key] = existingMetric
	storage.Mu.Unlock()

	repo := NewMetricMemoryGetRepository(storage)

	t.Run("existing metric", func(t *testing.T) {
		result, err := repo.Get(context.Background(), key)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingMetric.ID, result.ID)
		assert.Equal(t, existingMetric.MType, result.MType)
	})

	t.Run("non-existent metric", func(t *testing.T) {
		nonExistentKey := types.MetricID{ID: "not_exist", MType: "gauge"}
		result, err := repo.Get(context.Background(), nonExistentKey)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
