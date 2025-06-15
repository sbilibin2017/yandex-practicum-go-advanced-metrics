package main

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
)

func run(
	ctx context.Context,
	config *configs.AgentConfig,
	loggerInitializeFunc func(level string) error,
	newAgentFunc func(config *configs.AgentConfig) (func(ctx context.Context) error, error),
	newRunContextFunc func(ctx context.Context) (context.Context, context.CancelFunc),
	runWorkerFunc func(ctx context.Context, worker func(ctx context.Context) error) error,
) error {
	if err := loggerInitializeFunc(config.LogLevel); err != nil {
		return err
	}

	worker, err := newAgentFunc(config)
	if err != nil {
		return err
	}

	ctx, cancel := newRunContextFunc(ctx)
	defer cancel()

	return runWorkerFunc(ctx, worker)
}
