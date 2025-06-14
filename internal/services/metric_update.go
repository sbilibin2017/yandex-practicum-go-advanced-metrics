package services

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricUpdateSaver defines the interface for saving metric updates.
type MetricUpdateSaver interface {
	// Send saves a single metric.
	Send(ctx context.Context, metrics types.Metrics) error
	// Save finalizes the saving process.
	Save(ctx context.Context) error
}

// MetricUpdateFilterer defines the interface for filtering metrics by IDs.
type MetricUpdateFilterer interface {
	// Filter returns a map of existing metrics filtered by given MetricIDs.
	Filter(ctx context.Context, ids []types.MetricID) (map[types.MetricID]types.Metrics, error)
}

// MetricUpdateService provides methods to update metrics,
// combining filtering existing metrics and saving updates.
type MetricUpdateService struct {
	saver    MetricUpdateSaver
	filterer MetricUpdateFilterer
}

// NewMetricUpdateService creates a new MetricUpdateService instance
// with the provided saver and filterer implementations.
func NewMetricUpdateService(
	saver MetricUpdateSaver,
	filterer MetricUpdateFilterer,
) *MetricUpdateService {
	return &MetricUpdateService{saver: saver, filterer: filterer}
}

// Update processes a batch of metric updates.
// It accumulates metrics, filters existing metrics by IDs,
// merges counter metrics, and saves the updated metrics.
func (svc *MetricUpdateService) Update(
	ctx context.Context,
	metrics []types.Metrics,
) error {
	metrics = accumulateMetrics(metrics)

	metricsMap, err := filterMetricsByIDs(ctx, svc.filterer, metrics)
	if err != nil {
		logger.Log.Errorw("Failed to filter metrics by ID", "error", err)
		return err
	}

	metrics = updateMetrics(metrics, metricsMap)

	err = saveMetrics(ctx, svc.saver, metrics)
	if err != nil {
		logger.Log.Errorw("Failed to save metrics", "error", err)
		return err
	}

	return nil
}

// updateMetrics merges new metrics with existing metrics from the map.
// For counter metrics, it accumulates the delta values.
func updateMetrics(
	metrics []types.Metrics,
	metricsMap map[types.MetricID]types.Metrics,
) []types.Metrics {
	updated := make([]types.Metrics, 0, len(metrics))
	for _, m := range metrics {
		if existing, found := metricsMap[types.MetricID{ID: m.ID, MType: m.MType}]; found {
			if m.MType == types.Counter && m.Delta != nil && existing.Delta != nil {
				*m.Delta += *existing.Delta
			}
		}
		updated = append(updated, m)
	}
	return updated
}

// saveMetrics sends each metric to the saver and finalizes the saving process.
func saveMetrics(ctx context.Context, s MetricUpdateSaver, metrics []types.Metrics) error {
	for _, m := range metrics {
		if err := s.Send(ctx, m); err != nil {
			logger.Log.Errorw("Failed to send metric", "id", m.ID, "error", err)
			return err
		}
	}
	if err := s.Save(ctx); err != nil {
		logger.Log.Errorw("Failed to finalize save", "error", err)
		return err
	}
	return nil
}

// filterMetricsByIDs collects the IDs from the metrics and
// returns the existing metrics from the filterer.
func filterMetricsByIDs(
	ctx context.Context,
	f MetricUpdateFilterer,
	metrics []types.Metrics,
) (map[types.MetricID]types.Metrics, error) {
	var ids []types.MetricID
	for _, m := range metrics {
		ids = append(ids, types.MetricID{ID: m.ID, MType: m.MType})
	}
	return f.Filter(ctx, ids)
}

// accumulateMetrics aggregates counter metrics by summing their Delta values,
// and leaves gauge metrics as is.
func accumulateMetrics(metrics []types.Metrics) []types.Metrics {
	accumulated := make(map[types.MetricID]types.Metrics)
	for _, m := range metrics {
		id := types.MetricID{ID: m.ID, MType: m.MType}
		if m.MType == types.Counter {
			if existing, ok := accumulated[id]; ok && existing.Delta != nil && m.Delta != nil {
				sum := *existing.Delta + *m.Delta
				existing.Delta = &sum
				accumulated[id] = existing
			} else {
				accumulated[id] = m
			}
		} else {
			accumulated[id] = m
		}
	}
	result := make([]types.Metrics, 0, len(accumulated))
	for _, v := range accumulated {
		result = append(result, v)
	}
	return result
}
