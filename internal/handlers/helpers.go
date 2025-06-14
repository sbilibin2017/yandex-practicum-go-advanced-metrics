package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

var (
	errInternalServerError = errors.New("internal server error")
)

func handleError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

func handleInternalServerError(w http.ResponseWriter) {
	handleError(w, errInternalServerError.Error(), http.StatusInternalServerError)
}
