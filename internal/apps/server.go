package apps

import (
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/handlers"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/repositories"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/routers"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/services"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/validators"
)

// NewServerApp initializes and returns a new HTTP server configured
// with metric update routes, repositories, services, and handlers.
//
// It accepts a ServerConfig pointer that provides server address and log level settings.
//
// The returned *http.Server is ready to listen and serve on the configured address,
// and uses a router that handles metric update HTTP POST requests.
func NewServerApp(config *configs.ServerConfig) (*http.Server, error) {
	data := make(map[types.MetricID]types.Metrics)

	metricMemoryFiltererRepository := repositories.NewMetricMemoryFiltererRepository(data)
	metricMemorySaverRepository := repositories.NewMetricMemorySaverRepository(data)

	metricUpdateService := services.NewMetricUpdateService(
		metricMemorySaverRepository,
		metricMemoryFiltererRepository,
	)

	metricUpdatePathHandler := handlers.NewMetricUpdatePathHandler(
		validators.ValidateMetricPath,
		validators.HandleMetricsValidationError,
		metricUpdateService,
	)

	metricsRouter := routers.NewMetricsRouter(metricUpdatePathHandler)

	srv := &http.Server{
		Addr:    config.Address,
		Handler: metricsRouter,
	}

	return srv, nil
}
