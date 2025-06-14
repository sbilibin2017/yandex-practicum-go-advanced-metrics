package repositories

import (
	"context"
	"sync"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricMemorySaverRepository is an in-memory implementation of a repository
// that stores metrics with thread-safe access using a read-write mutex.
type MetricMemorySaverRepository struct {
	mu   sync.RWMutex
	data map[types.MetricID]types.Metrics
}

// NewMetricMemorySaverRepository creates a new MetricMemorySaverRepository
// initialized with the provided data map.
//
// Parameters:
//   - data: a map from MetricID to Metrics to initialize the repository.
//
// Returns:
//   - pointer to a new MetricMemorySaverRepository instance.
func NewMetricMemorySaverRepository(data map[types.MetricID]types.Metrics) *MetricMemorySaverRepository {
	return &MetricMemorySaverRepository{
		data: data,
		mu:   sync.RWMutex{},
	}
}

// Send stores or updates a metric in the repository.
//
// It locks the repository for writing, inserts or updates the metric in the internal map,
// then releases the lock.
//
// Parameters:
//   - ctx: context for request scoping (not used here but included for interface compatibility).
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
//   - error: always nil in this implementation.
func (repo *MetricMemorySaverRepository) Save(ctx context.Context) error {
	return nil
}
