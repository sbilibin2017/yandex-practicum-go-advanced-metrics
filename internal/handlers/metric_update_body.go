package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricBodyUpdater interface {
	Update(ctx context.Context, metrics []types.Metrics) error
}

func NewMetricUpdateBodyHandler(
	valFunc func(metric types.Metrics) error,
	errValHandlerFunc func(err error) *types.APIError,
	svc MetricBodyUpdater,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		var metric types.Metrics
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&metric); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err := valFunc(metric)
		if apiErr := errValHandlerFunc(err); apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
			return
		}

		if err := svc.Update(r.Context(), []types.Metrics{metric}); err != nil {
			handleInternalServerError(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(metric); err != nil {
			handleInternalServerError(w)
			return
		}
	}
}
