package main

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/apps"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/runners"
)

func main() {
	config, err := parseFlags()
	if err != nil {
		panic(err)
	}

	err = run(
		context.Background(),
		config,
		logger.Initialize,
		apps.NewAgentApp,
		runners.NewRunContext,
		runners.RunWorker,
	)
	if err != nil {
		panic(err)
	}
}
