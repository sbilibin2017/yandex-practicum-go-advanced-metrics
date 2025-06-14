package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/runners"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	mockServer := &http.Server{}

	tests := []struct {
		name            string
		loggerErr       error
		serverAppErr    error
		runServerErr    error
		wantErr         bool
		wantErrContains string
	}{
		{
			name:            "logger initialization fails",
			loggerErr:       errors.New("logger error"),
			wantErr:         true,
			wantErrContains: "logger error",
		},
		{
			name:            "new server app fails",
			loggerErr:       nil,
			serverAppErr:    errors.New("server app error"),
			wantErr:         true,
			wantErrContains: "server app error",
		},
		{
			name:            "run server fails",
			loggerErr:       nil,
			serverAppErr:    nil,
			runServerErr:    errors.New("run server error"),
			wantErr:         true,
			wantErrContains: "run server error",
		},
		{
			name:         "success case",
			loggerErr:    nil,
			serverAppErr: nil,
			runServerErr: nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			loggerInitializeFunc := func(level string) error {
				return tt.loggerErr
			}
			newServerFunc := func(cfg *configs.ServerConfig) (*http.Server, error) {
				if tt.serverAppErr != nil {
					return nil, tt.serverAppErr
				}
				return mockServer, nil
			}
			newRunContextFunc := func(ctx context.Context) (context.Context, context.CancelFunc) {
				return context.WithCancel(ctx)
			}
			runServerFunc := func(ctx context.Context, srv runners.Server) error {
				return tt.runServerErr
			}

			err := run(
				context.Background(),
				&configs.ServerConfig{LogLevel: "info"},
				loggerInitializeFunc,
				newServerFunc,
				newRunContextFunc,
				runServerFunc,
			)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
