package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricMemoryGetRepository struct {
	storage *engines.MemoryStorage[types.MetricID, types.Metrics]
}

func NewMetricMemoryGetRepository(
	storage *engines.MemoryStorage[types.MetricID, types.Metrics],
) *MetricMemoryGetRepository {
	return &MetricMemoryGetRepository{
		storage: storage,
	}
}

func (repo *MetricMemoryGetRepository) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	repo.storage.Mu.RLock()
	defer repo.storage.Mu.RUnlock()

	metric, ok := repo.storage.Data[id]
	if !ok {
		return nil, nil
	}

	return &metric, nil
}
