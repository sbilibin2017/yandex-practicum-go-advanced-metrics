package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricMemorySaveRepository struct {
	storage *engines.MemoryStorage[types.MetricID, types.Metrics]
}

func NewMetricMemorySaveRepository(
	storage *engines.MemoryStorage[types.MetricID, types.Metrics],
) *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{
		storage: storage,
	}
}

func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context,
	metrics types.Metrics,
) error {
	repo.storage.Mu.Lock()
	defer repo.storage.Mu.Unlock()

	repo.storage.Data[types.MetricID{
		ID:    metrics.ID,
		MType: metrics.MType,
	}] = metrics

	return nil
}
