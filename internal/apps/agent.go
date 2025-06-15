package apps

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/facades"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/workers"
)

func NewAgentApp(config *configs.AgentConfig) (func(ctx context.Context) error, error) {
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
