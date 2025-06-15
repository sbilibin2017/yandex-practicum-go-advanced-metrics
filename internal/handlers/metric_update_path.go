package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricPathUpdater interface {
	Update(ctx context.Context, metrics []types.Metrics) error
}

func NewMetricUpdatePathHandler(
	valFunc func(metricName string, metricType string, metricValue string) error,
	errValHandlerFunc func(err error) *types.APIError,
	svc MetricPathUpdater,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := getURLParam(r, "type")
		metricName := getURLParam(r, "name")
		metricValue := getURLParam(r, "value")

		err := valFunc(metricName, metricType, metricValue)

		apiErr := errValHandlerFunc(err)
		if apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
			return
		}

		metric := newMetrics(metricType, metricName, metricValue)

		if err := svc.Update(r.Context(), []types.Metrics{*metric}); err != nil {
			handleInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func newMetrics(metricType, metricName, metricValue string) *types.Metrics {
	m := &types.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case types.Counter:
		if delta, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
			m.Delta = &delta
		}
	case types.Gauge:
		if value, err := strconv.ParseFloat(metricValue, 64); err == nil {
			m.Value = &value
		}
	}
	return m
}
