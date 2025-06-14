package repositories

import (
	"context"
	"sync"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricMemoryFiltererRepository is an in-memory repository implementation
// that supports filtering metrics by their IDs with thread-safe access,
// sharing the lock externally.
type MetricMemoryFiltererRepository struct {
	mu   *sync.RWMutex
	data map[types.MetricID]types.Metrics
}

// NewMetricMemoryFiltererRepository creates a new MetricMemoryFiltererRepository
// initialized with the provided data map and shared mutex.
//
// Parameters:
//   - data: map from MetricID to Metrics to initialize the repository.
//   - mu: pointer to shared RWMutex for synchronization.
//
// Returns:
//   - pointer to a new MetricMemoryFiltererRepository instance.
func NewMetricMemoryFiltererRepository() *MetricMemoryFiltererRepository {
	return &MetricMemoryFiltererRepository{
		data: data,
		mu:   mu,
	}
}

// Filter returns a map of Metrics filtered by the given list of MetricIDs.
//
// It acquires a read lock on the shared mutex to safely access the internal data.
//
// Parameters:
//   - ctx: context for request scoping (unused).
//   - ids: slice of MetricID to filter metrics.
//
// Returns:
//   - map of MetricID to Metrics for the requested IDs found in the repository.
//   - error, always nil.
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
