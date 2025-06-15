package runners

import (
	"context"
	"time"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

const defaultShutdownTimeout = 5 * time.Second

func RunServer(ctx context.Context, srv Server) error {
	errCh := make(chan error, 1)

	logger.Log.Infow("Starting server")

	go func() {
		errCh <- srv.ListenAndServe()
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
		return nil

	case err := <-errCh:
		if err != nil {
			logger.Log.Errorw("Server stopped with error", "error", err)
		}
		return err
	}
}
