package runners

import (
	"context"
	"errors"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
)

// RunWorker runs a given worker function in a separate goroutine,
// monitors its execution, and handles panics gracefully.
//
// It listens for cancellation of the provided context, and returns the context's error
// if cancelled before the worker finishes.
//
// If the worker panics, the panic is recovered, logged, and an error is returned.
//
// The function returns nil if the worker completes successfully without panicking.
//
// Logging is performed for worker start, panic recovery, normal completion, and context cancellation.
func RunWorker(ctx context.Context, worker func(ctx context.Context)) error {
	errCh := make(chan error, 1)

	logger.Log.Infow("Starting worker")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Errorw("Worker panicked", "recover", r)
				errCh <- errors.New("worker panicked")
			}
		}()

		worker(ctx)
		logger.Log.Infow("Worker finished without panic")
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		logger.Log.Infow("Worker context cancelled", "reason", ctx.Err())
		return ctx.Err()

	case err := <-errCh:
		if err != nil {
			logger.Log.Errorw("Worker exited with error", "error", err)
		} else {
			logger.Log.Infow("Worker exited cleanly")
		}
		return err
	}
}
