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
		worker     func(ctx context.Context) error
		ctxTimeout time.Duration
		wantErr    error
	}{
		{
			name: "worker completes successfully",
			worker: func(ctx context.Context) error {
				return nil
			},
			ctxTimeout: 100 * time.Millisecond,
			wantErr:    nil,
		},
		{
			name: "worker returns after delay",
			worker: func(ctx context.Context) error {
				time.Sleep(50 * time.Millisecond)
				return nil
			},
			ctxTimeout: 200 * time.Millisecond,
			wantErr:    nil,
		},
		{
			name: "worker returns error",
			worker: func(ctx context.Context) error {
				return errors.New("worker error")
			},
			ctxTimeout: 100 * time.Millisecond,
			wantErr:    errors.New("worker error"),
		},
		{
			name: "context canceled before worker finishes",
			worker: func(ctx context.Context) error {
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			ctxTimeout: 50 * time.Millisecond,
			wantErr:    nil, // RunWorker returns nil on ctx.Done()
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
			} else {
				require.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
