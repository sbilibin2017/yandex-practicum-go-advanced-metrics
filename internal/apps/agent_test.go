package apps

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/stretchr/testify/require"
)

func TestNewAgentApp(t *testing.T) {
	config := &configs.AgentConfig{
		ServerAddress:  "http://localhost:8080",
		ServerEndpoint: "/update",
		PollInterval:   1,
		ReportInterval: 1,
		NumWorkers:     1,
	}

	workerFunc, err := NewAgentApp(config)
	require.NoError(t, err)
	require.NotNil(t, workerFunc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the worker function in a goroutine to check it does not panic or block indefinitely
	go workerFunc(ctx)
}
