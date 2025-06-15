package services

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricListLister interface {
	List(ctx context.Context) ([]types.Metrics, error)
}

type MetricListService struct {
	lister MetricListLister
}

func NewMetricListService(
	lister MetricListLister,
) *MetricListService {
	return &MetricListService{lister: lister}
}

func (svc *MetricListService) List(
	ctx context.Context,
) ([]types.Metrics, error) {
	logger.Log.Debug("MetricListService.List called")

	metrics, err := svc.lister.List(ctx)
	if err != nil {
		logger.Log.Errorw("Failed to list metrics",
			"error", err,
		)
		return nil, err
	}

	logger.Log.Debugw("Metrics listed successfully",
		"count", len(metrics),
	)

	return metrics, nil
}
