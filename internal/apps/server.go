package apps

import (
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/engines"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/handlers"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/middlewares"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/repositories"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/routers"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/services"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/validators"
)

func NewServerApp(config *configs.ServerConfig) (*http.Server, error) {
	memStorage := engines.NewMemoryStorage[types.MetricID, types.Metrics]()

	metricMemoryGetRepository := repositories.NewMetricMemoryGetRepository(memStorage)
	metricMemorySaverRepository := repositories.NewMetricMemorySaveRepository(memStorage)
	metricMemoryListerRepository := repositories.NewMetricMemoryListRepository(memStorage)

	metricUpdateService := services.NewMetricUpdateService(
		metricMemorySaverRepository,
		metricMemoryGetRepository,
	)
	metricGetService := services.NewMetricGetService(
		metricMemoryGetRepository,
	)
	metricListService := services.NewMetricListService(
		metricMemoryListerRepository,
	)

	metricUpdatePathHandler := handlers.NewMetricUpdatePathHandler(
		validators.ValidateMetricPath,
		validators.HandleMetricsValidationError,
		metricUpdateService,
	)
	metricGetPathHandler := handlers.NewMetricGetPathHandler(
		validators.ValidateMetricIDPath,
		validators.HandleMetricsValidationError,
		metricGetService,
	)
	metricListHTMLHandler := handlers.NewMetricListHTMLHandler(
		validators.HandleMetricsValidationError,
		metricListService,
	)

	middlewares := []func(next http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
	}

	metricsRouter := routers.NewMetricsRouter(
		metricUpdatePathHandler,
		metricGetPathHandler,
		metricListHTMLHandler,
		middlewares...,
	)

	srv := &http.Server{
		Addr:    config.Address,
		Handler: metricsRouter,
	}

	return srv, nil
}
