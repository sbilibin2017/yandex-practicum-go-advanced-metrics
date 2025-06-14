package apps

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/facades"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/workers"
)

// NewAgentApp creates and returns a metric agent worker function along with an error.
//
// It sets up an HTTP client, constructs a MetricUpdateFacade using the client and configuration,
// and initializes a metric agent worker with the provided polling interval, reporting interval,
// and number of concurrent workers.
//
// Parameters:
//   - config: Agent configuration containing server address, endpoint, polling/reporting intervals, and number of workers.
//
// Returns:
//   - func(ctx context.Context): A function that starts the metric agent worker when called with a context.
//   - error: Always returns nil in the current implementation but reserved for future error handling.
func NewAgentApp(config *configs.AgentConfig) (func(ctx context.Context), error) {
	client := resty.New()
	metricUpdateFacade := facades.NewMetricUpdateFacade(client, config.ServerAddress, config.ServerEndpoint)
	worker := workers.NewMetricAgentWorker(
		metricUpdateFacade,
		config.PollInterval,
		config.ReportInterval,
		config.NumWorkers,
	)
	return worker, nil
}
