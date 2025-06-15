package repositories

import (
	"context"
	"sort"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricMemoryListRepository struct {
	storage *engines.MemoryStorage[types.MetricID, types.Metrics]
}

func NewMetricMemoryListRepository(
	storage *engines.MemoryStorage[types.MetricID, types.Metrics],
) *MetricMemoryListRepository {
	return &MetricMemoryListRepository{
		storage: storage,
	}
}

func (repo *MetricMemoryListRepository) List(
	ctx context.Context,
) ([]types.Metrics, error) {
	repo.storage.Mu.RLock()
	defer repo.storage.Mu.RUnlock()

	metrics := make([]types.Metrics, 0, len(repo.storage.Data))
	for _, metric := range repo.storage.Data {
		metrics = append(metrics, metric)
	}

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ID < metrics[j].ID
	})

	return metrics, nil
}
