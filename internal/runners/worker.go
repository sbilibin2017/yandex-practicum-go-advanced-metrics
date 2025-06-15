package runners

import (
	"context"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
)

func RunWorker(ctx context.Context, worker func(ctx context.Context) error) error {
	errCh := make(chan error, 1)

	logger.Log.Infow("Starting worker")

	go func() {
		errCh <- worker(ctx)
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info("Worker stopped successfully")
		return nil
	case err := <-errCh:
		if err != nil {
			logger.Log.Errorw("Worker stopped with error", "error", err)
		}
		return err
	}
}
