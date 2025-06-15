package runners

import (
	"context"
	"os/signal"
	"syscall"
)

func NewRunContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	return ctx, stop
}
