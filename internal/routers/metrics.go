package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMetricsRouter(
	metricUpdatePathHandler http.HandlerFunc,
	metricValuePathHandler http.HandlerFunc,
	metricsListHandler http.HandlerFunc, // ← Новый параметр
	middlewares ...func(http.Handler) http.Handler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares...)

	r.Post("/update/{type}/{name}/{value}", metricUpdatePathHandler)
	r.Post("/update/{type}/{name}", metricUpdatePathHandler)

	r.Get("/value/{type}/{name}", metricValuePathHandler)
	r.Get("/value/{type}", metricValuePathHandler)

	r.Get("/", metricsListHandler)

	return r
}
