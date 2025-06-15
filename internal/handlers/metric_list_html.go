package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricHTMLLister interface {
	List(ctx context.Context) ([]types.Metrics, error)
}

func NewMetricListHTMLHandler(
	errHandlerFunc func(err error) *types.APIError,
	svc MetricHTMLLister,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := svc.List(r.Context())

		apiErr := errHandlerFunc(err)
		if apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(types.GetMetricsHTML(metrics)))
	}
}
