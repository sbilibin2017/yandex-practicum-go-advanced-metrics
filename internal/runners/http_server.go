package runners

import (
	"context"
	"errors"
	"time"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

const defaultShutdownTimeout = 5 * time.Second

// RunServer starts the given Server in a separate goroutine and manages its lifecycle.
//
// It listens for the provided context to be cancelled, and upon cancellation,
// it initiates a graceful shutdown of the server with a timeout of 5 seconds.
//
// The function returns an error if the server fails to start, stops with an error,
// or if the shutdown process encounters an error. If the server stops gracefully,
// or the context is cancelled and the server shuts down cleanly, it returns nil or
// the context cancellation error.
//
// The server's ListenAndServe method is expected to return either nil or an error
// that is not context.Canceled. If the error is context.Canceled, it is treated as a
// graceful shutdown.
//
// This function also logs server lifecycle events.
func RunServer(ctx context.Context, srv Server) error {
	errCh := make(chan error, 1)

	logger.Log.Infow("Starting server")

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Log.Errorw("Server exited with error", "error", err)
			errCh <- err
		} else {
			logger.Log.Infow("Server shutdown gracefully")
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Infow("Context cancelled, initiating shutdown")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Log.Errorw("Server shutdown error", "error", err)
			return err
		}

		logger.Log.Infow("Server shutdown complete")
		return ctx.Err()

	case err := <-errCh:
		if err != nil {
			logger.Log.Errorw("Server stopped with error", "error", err)
		} else {
			logger.Log.Infow("Server stopped cleanly")
		}
		return err
	}
}
