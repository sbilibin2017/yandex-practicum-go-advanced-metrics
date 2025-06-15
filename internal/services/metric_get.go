package services

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricGetGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

type MetricGetService struct {
	getter MetricGetGetter
}

func NewMetricGetService(
	getter MetricGetGetter,
) *MetricGetService {
	return &MetricGetService{getter: getter}
}

func (svc *MetricGetService) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	logger.Log.Infow("MetricGetService.Get called",
		"id", id.ID,
		"type", id.MType,
	)

	metric, err := svc.getter.Get(ctx, id)
	if err != nil {
		logger.Log.Errorw("Failed to get metric",
			"id", id.ID,
			"type", id.MType,
			"error", err,
		)
		return nil, err
	}

	if metric == nil {
		logger.Log.Warnw("Metric not found",
			"id", id.ID,
			"type", id.MType,
		)
		return nil, errors.ErrMetricNotFound
	}

	logger.Log.Infow("Metric retrieved successfully",
		"id", metric.ID,
		"type", metric.MType,
	)

	return metric, nil
}
