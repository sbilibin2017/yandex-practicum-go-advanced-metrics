package services

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricUpdateSaver interface {
	Save(ctx context.Context, metrics types.Metrics) error
}

type MetricUpdateGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

type MetricUpdateService struct {
	saver  MetricUpdateSaver
	getter MetricUpdateGetter
}

func NewMetricUpdateService(
	saver MetricUpdateSaver,
	getter MetricUpdateGetter,
) *MetricUpdateService {
	return &MetricUpdateService{saver: saver, getter: getter}
}

func (svc *MetricUpdateService) Update(
	ctx context.Context,
	metrics []types.Metrics,
) error {
	for _, m := range metrics {
		if m.MType == types.Counter {
			existing, err := svc.getter.Get(ctx, types.MetricID{ID: m.ID, MType: m.MType})
			if err != nil {
				logger.Log.Errorw("Failed to get existing metric",
					"id", m.ID,
					"error", err,
				)
				return err
			}

			if existing != nil && m.Delta != nil && existing.Delta != nil {
				*m.Delta += *existing.Delta
			}
		}

		if err := svc.saver.Save(ctx, m); err != nil {
			logger.Log.Errorw("Failed to save metric",
				"id", m.ID,
				"error", err,
			)
			return err
		}
	}

	return nil
}
