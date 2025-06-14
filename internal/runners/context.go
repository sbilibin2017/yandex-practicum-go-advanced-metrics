package runners

import (
	"context"
	"os/signal"
	"syscall"
)

// NewRunContext returns a context that is canceled when one of the specified
// OS signals (SIGINT, SIGTERM, SIGQUIT) is received.
// It wraps the provided parent context and returns the derived context
// along with a cancel function to stop signal notification.
//
// This is useful for gracefully shutting down applications on termination signals.
func NewRunContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	return ctx, stop
}
