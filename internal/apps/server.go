package apps

import (
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/handlers"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/repositories"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/routers"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/services"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/validators"
)

// NewServerApp initializes and returns a new HTTP server configured to handle metric update requests.
//
// It sets up in-memory repositories for metric data storage and filtering,
// creates the metric update service layer, and registers HTTP handlers and routers.
//
// The server listens on the address specified in the provided ServerConfig.
//
// Parameters:
//   - config: Pointer to ServerConfig containing server address and logging settings.
//
// Returns:
//   - *http.Server: Configured HTTP server instance ready to listen and serve requests.
//   - error: An error if the server initialization fails (currently always nil).
func NewServerApp(config *configs.ServerConfig) (*http.Server, error) {

	metricMemoryFiltererRepository := repositories.NewMetricMemoryFiltererRepository()
	metricMemorySaverRepository := repositories.NewMetricMemorySaverRepository()

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
