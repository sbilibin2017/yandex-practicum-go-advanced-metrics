package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricBodyGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

func NewMetricGetBodyHandler(
	valFunc func(id types.MetricID) error,
	errValHandlerFunc func(err error) *types.APIError,
	svc MetricBodyGetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		var id types.MetricID
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&id); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		err := valFunc(id)
		if apiErr := errValHandlerFunc(err); apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
			return
		}

		metric, err := svc.Get(r.Context(), id)
		if apiErr := errValHandlerFunc(err); apiErr != nil {
			handleError(w, apiErr.Message, apiErr.Code)
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
