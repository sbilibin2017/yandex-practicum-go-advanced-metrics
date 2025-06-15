package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"
)

func getURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func handleError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

func handleInternalServerError(w http.ResponseWriter) {
	handleError(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
}
