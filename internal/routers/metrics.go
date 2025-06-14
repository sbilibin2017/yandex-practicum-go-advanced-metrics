package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewMetricsRouter creates and returns a new chi router configured
// with routes for updating metrics via HTTP POST requests.
//
// It registers the given metricUpdatePathHandler to handle POST requests
// at the following endpoints:
//   - /update/{type}/{name}/{value}
//   - /update/{type}/{name}
//
// The function accepts optional middleware handlers which are applied
// to all routes in the router.
//
// Parameters:
//   - metricUpdatePathHandler: the HTTP handler function responsible for
//     processing metric update requests.
//   - middlewares: variadic middleware functions to be applied to the router.
//
// Returns:
//   - *chi.Mux: the configured router instance ready to be used in an HTTP server.
func NewMetricsRouter(
	metricUpdatePathHandler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Post("/update/{type}/{name}/{value}", metricUpdatePathHandler)
	r.Post("/update/{type}/{name}", metricUpdatePathHandler)
	return r
}
