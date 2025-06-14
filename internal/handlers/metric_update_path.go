package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricPathUpdater defines an interface for updating metrics based on path parameters.
type MetricPathUpdater interface {
	// Update processes and saves the provided metrics.
	Update(ctx context.Context, metrics []types.Metrics) error
}

// NewMetricUpdatePathHandler returns an HTTP handler that processes metric updates
// based on URL path parameters.
//
// valFunc validates metric parameters: metric name, type, and value.
//
// errValHandlerFunc converts validation errors into APIError responses.
//
// svc is a service implementing MetricPathUpdater that performs the update.
//
// The handler extracts "type", "name", and "value" parameters from the URL,
// validates them using valFunc, converts validation errors using errValHandlerFunc,
// and, if valid, creates a Metrics struct and calls svc.Update.
// On success, it responds with HTTP 200 OK; otherwise, appropriate error responses are returned.
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

// newMetrics creates a Metrics struct from the given metric type, name, and string value.
// It parses the value into the appropriate numeric type depending on metricType (Counter or Gauge).
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
