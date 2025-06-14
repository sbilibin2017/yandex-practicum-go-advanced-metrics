package main

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/runners"
)

// run sets up and runs the HTTP server with the provided configuration and functions.
//
// It initializes the logger using loggerInitializeFunc with the log level from config.
// Then it creates a new HTTP server using newServerFunc.
// It creates a cancellable context for server execution using newRunContextFunc.
// Finally, it runs the server with runServerFunc, blocking until the server stops or the context is cancelled.
//
// Parameters:
// - ctx: the base context for running the server.
// - config: server configuration containing address and log level.
// - loggerInitializeFunc: function to initialize the logger with a given log level.
// - newServerFunc: function to create a new HTTP server instance.
// - newRunContextFunc: function to create a cancellable context that listens for OS signals.
// - runServerFunc: function that runs the server and handles shutdown.
//
// Returns:
// - error if any step fails, otherwise nil after server shutdown.
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
