package main

import (
	"context"
	"errors"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ctxKey string

const myKey ctxKey = "key"

func TestRun_Success(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.AgentConfig{LogLevel: "debug"}

	loggerCalled := false
	loggerInitializeFunc := func(level string) error {
		loggerCalled = true
		assert.Equal(t, "debug", level)
		return nil
	}

	workerFunc := func(ctx context.Context) error {
		return nil
	}

	newAgentFunc := func(c *configs.AgentConfig) (func(ctx context.Context) error, error) {
		assert.Equal(t, cfg, c)
		return workerFunc, nil
	}

	runCtx := context.WithValue(context.Background(), myKey, "value")
	cancelCalled := false
	newRunContextFunc := func(ctx context.Context) (context.Context, context.CancelFunc) {
		return runCtx, func() { cancelCalled = true }
	}

	runWorkerCalled := false
	runWorkerFunc := func(ctx context.Context, worker func(ctx context.Context) error) error {
		runWorkerCalled = true
		assert.Equal(t, runCtx, ctx)
		assert.NotNil(t, worker)
		return nil
	}

	err := run(ctx, cfg, loggerInitializeFunc, newAgentFunc, newRunContextFunc, runWorkerFunc)
	require.NoError(t, err)

	assert.True(t, loggerCalled, "loggerInitializeFunc should be called")
	assert.True(t, cancelCalled, "cancel func should be called")
	assert.True(t, runWorkerCalled, "runWorkerFunc should be called")
}

func TestRun_LoggerError(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.AgentConfig{LogLevel: "error"}

	wantErr := errors.New("logger init failed")
	loggerInitializeFunc := func(level string) error {
		assert.Equal(t, "error", level)
		return wantErr
	}

	newAgentFunc := func(c *configs.AgentConfig) (func(ctx context.Context) error, error) {
		t.Fatal("newAgentFunc should not be called if logger fails")
		return nil, nil
	}

	newRunContextFunc := func(ctx context.Context) (context.Context, context.CancelFunc) {
		t.Fatal("newRunContextFunc should not be called if logger fails")
		return ctx, func() {}
	}

	runWorkerFunc := func(ctx context.Context, worker func(ctx context.Context) error) error {
		t.Fatal("runWorkerFunc should not be called if logger fails")
		return nil
	}

	err := run(ctx, cfg, loggerInitializeFunc, newAgentFunc, newRunContextFunc, runWorkerFunc)
	assert.ErrorIs(t, err, wantErr)
}

func TestRun_NewAgentFuncError(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.AgentConfig{LogLevel: "info"}

	wantErr := errors.New("failed to create agent")

	loggerInitializeFunc := func(level string) error {
		assert.Equal(t, "info", level)
		return nil
	}

	newAgentFunc := func(c *configs.AgentConfig) (func(ctx context.Context) error, error) {
		assert.Equal(t, cfg, c)
		return nil, wantErr
	}

	newRunContextFunc := func(ctx context.Context) (context.Context, context.CancelFunc) {
		return ctx, func() {}
	}

	runWorkerFunc := func(ctx context.Context, worker func(ctx context.Context) error) error {
		t.Fatal("runWorkerFunc should not be called when newAgentFunc fails")
		return nil
	}

	err := run(ctx, cfg, loggerInitializeFunc, newAgentFunc, newRunContextFunc, runWorkerFunc)
	assert.ErrorIs(t, err, wantErr)
}
