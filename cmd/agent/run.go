package main

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
)

// run initializes and runs the metric agent lifecycle.
//
// It performs the following steps:
//  1. Initializes the logger using loggerInitializeFunc with the log level from the config.
//  2. Creates the agent worker function by calling newAgentFunc with the given config.
//  3. Creates a cancellable context by calling newRunContextFunc, typically listening for OS signals.
//  4. Runs the worker function via runWorkerFunc using the cancellable context.
//
// Parameters:
//   - ctx: The base context for running the agent.
//   - config: Configuration settings for the agent.
//   - loggerInitializeFunc: Function to initialize logging, accepts log level string.
//   - newAgentFunc: Function that returns the agent worker function given the config.
//   - newRunContextFunc: Function to create a cancellable context from the base context.
//   - runWorkerFunc: Function to run the worker with the provided context.
//
// Returns:
//   - error: An error returned by any of the initialization or run steps, or nil on success.
func run(
	ctx context.Context,
	config *configs.AgentConfig,
	loggerInitializeFunc func(level string) error,
	newAgentFunc func(config *configs.AgentConfig) (func(ctx context.Context), error),
	newRunContextFunc func(ctx context.Context) (context.Context, context.CancelFunc),
	runWorkerFunc func(ctx context.Context, worker func(ctx context.Context)) error,
) error {
	// Initialize logger with config log level
	if err := loggerInitializeFunc(config.LogLevel); err != nil {
		return err
	}

	// Create the agent worker function
	worker, err := newAgentFunc(config)
	if err != nil {
		return err
	}

	// Create cancellable context (e.g., listens for OS signals)
	ctx, cancel := newRunContextFunc(ctx)
	defer cancel()

	// Run the worker with the context, block until done or error
	return runWorkerFunc(ctx, worker)
}
