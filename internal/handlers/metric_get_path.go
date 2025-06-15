package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricPathGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

func NewMetricGetPathHandler(
	valFunc func(metricName string, metricType string) error,
	errHandlerFunc func(err error) *types.APIError,
	svc MetricPathGetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := getURLParam(r, "type")
		metricName := getURLParam(r, "name")

		err := valFunc(metricName, metricType)

		apiErr := errHandlerFunc(err)
		if apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
			return
		}

		metricID := newMetricID(metricType, metricName)

		metrics, err := svc.Get(r.Context(), *metricID)

		apiErr = errHandlerFunc(err)
		if apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(types.GetMetricStringValue(metrics)))
	}
}

func newMetricID(metricType, metricName string) *types.MetricID {
	return &types.MetricID{
		ID:    metricName,
		MType: metricType,
	}
}
