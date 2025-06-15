package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMetricsRouter(
	metricUpdatePathHandler http.HandlerFunc,
	metricUpdateBodyHandler http.HandlerFunc,
	metricValuePathHandler http.HandlerFunc,
	metricValueBodyHandler http.HandlerFunc,
	metricsListHandler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares...)

	r.Post("/update/{type}/{name}/{value}", metricUpdatePathHandler)
	r.Post("/update/{type}/{name}", metricUpdatePathHandler)
	r.Post("/update/", metricUpdateBodyHandler)

	r.Get("/value/{type}/{name}", metricValuePathHandler)
	r.Get("/value/{type}", metricValuePathHandler)
	r.Post("/value/", metricValueBodyHandler)

	r.Get("/", metricsListHandler)

	return r
}
