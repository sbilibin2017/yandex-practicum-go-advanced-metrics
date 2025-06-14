package runners

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunWorker(t *testing.T) {
	tests := []struct {
		name       string
		worker     func(ctx context.Context)
		ctxTimeout time.Duration
		wantErr    error
	}{
		{
			name: "worker completes successfully",
			worker: func(ctx context.Context) {
			},
			ctxTimeout: 100 * time.Millisecond,
			wantErr:    nil,
		},
		{
			name: "worker panics and is recovered",
			worker: func(ctx context.Context) {
				panic("unexpected panic")
			},
			ctxTimeout: 100 * time.Millisecond,
			wantErr:    errors.New("worker panicked"),
		},
		{
			name: "context canceled before worker finishes",
			worker: func(ctx context.Context) {
				<-ctx.Done()
			},
			ctxTimeout: 10 * time.Millisecond,
			wantErr:    context.Canceled, // logical expectation, but test for both below
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.ctxTimeout)
			defer cancel()

			err := RunWorker(ctx, tt.worker)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else if tt.name == "context canceled before worker finishes" {
				// Accept context.Canceled or context.DeadlineExceeded
				require.Error(t, err)
				require.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
					"expected context.Canceled or context.DeadlineExceeded, got %v", err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
