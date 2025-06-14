package repositories

import (
	"context"
	"sync"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricMemorySaverRepository stores metrics with thread-safe access,
// sharing the lock externally to coordinate concurrent access.
type MetricMemorySaverRepository struct {
	mu   *sync.RWMutex
	data map[types.MetricID]types.Metrics
}

// NewMetricMemorySaverRepository creates a new MetricMemorySaverRepository
// with the given data map and shared mutex.
//
// Parameters:
//   - data: map from MetricID to Metrics to initialize storage.
//   - mu: pointer to a shared RWMutex for synchronizing access.
//
// Returns:
//   - pointer to a new MetricMemorySaverRepository instance.
func NewMetricMemorySaverRepository() *MetricMemorySaverRepository {
	return &MetricMemorySaverRepository{
		data: data,
		mu:   mu,
	}
}

// Send stores or updates a metric in the repository.
//
// It locks the repository for writing using the shared mutex,
// inserts or updates the metric in the internal map,
// then releases the lock.
//
// Parameters:
//   - ctx: context for request scoping (unused).
//   - metrics: the Metrics object to store.
//
// Returns:
//   - error: always nil in this implementation.
func (repo *MetricMemorySaverRepository) Send(ctx context.Context, metrics types.Metrics) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.data[types.MetricID{ID: metrics.ID, MType: metrics.MType}] = metrics
	return nil
}

// Save is a no-op for this in-memory repository.
//
// Parameters:
//   - ctx: context for request scoping.
//
// Returns:
//   - error: always nil.
func (repo *MetricMemorySaverRepository) Save(ctx context.Context) error {
	return nil
}
