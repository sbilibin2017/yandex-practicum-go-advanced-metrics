package main

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/runners"
)

func run(
	ctx context.Context,
	config *configs.ServerConfig,
	loggerInitializeFunc func(level string) error,
	newServerFunc func(*configs.ServerConfig) (*http.Server, error),
	newRunContextFunc func(ctx context.Context) (context.Context, context.CancelFunc),
	runServerFunc func(ctx context.Context, srv runners.Server) error,
) error {
	err := loggerInitializeFunc(config.LogLevel)
	if err != nil {
		return err
	}

	srv, err := newServerFunc(config)
	if err != nil {
		return err
	}

	ctx, cancel := newRunContextFunc(ctx)
	defer cancel()

	return runServerFunc(ctx, srv)
}
