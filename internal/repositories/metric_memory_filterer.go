package repositories

import (
	"context"
	"sync"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricMemoryFiltererRepository is an in-memory repository implementation
// that supports filtering metrics by their IDs with thread-safe access.
type MetricMemoryFiltererRepository struct {
	mu   sync.RWMutex
	data map[types.MetricID]types.Metrics
}

// NewMetricMemoryFiltererRepository creates a new MetricMemoryFiltererRepository
// initialized with the provided data map.
//
// Parameters:
//   - data: a map from MetricID to Metrics to initialize the repository.
//
// Returns:
//   - pointer to a new MetricMemoryFiltererRepository instance.
func NewMetricMemoryFiltererRepository(data map[types.MetricID]types.Metrics) *MetricMemoryFiltererRepository {
	return &MetricMemoryFiltererRepository{
		data: data,
		mu:   sync.RWMutex{},
	}
}

// Filter returns a map of Metrics filtered by the given list of MetricIDs.
//
// It acquires a read lock on the repository to safely access the internal data.
//
// Parameters:
//   - ctx: context for request scoping (not used here but included for interface compatibility).
//   - ids: slice of MetricID to filter metrics.
//
// Returns:
//   - map of MetricID to Metrics for the requested IDs found in the repository.
//   - error, always nil in this implementation.
func (repo *MetricMemoryFiltererRepository) Filter(
	ctx context.Context,
	ids []types.MetricID,
) (map[types.MetricID]types.Metrics, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	result := make(map[types.MetricID]types.Metrics, len(ids))
	for _, id := range ids {
		if metric, found := repo.data[id]; found {
			result[id] = metric
		}
	}
	return result, nil
}
