package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/apps"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/runners"
)

type MockAgentSuite struct {
	suite.Suite

	serverURL string
	ts        *httptest.Server

	mu      sync.Mutex
	metrics []MockMetric
}

type MockMetric struct {
	Type  string
	Name  string
	Value string
}

func (s *MockAgentSuite) SetupSuite() {
	err := logger.Initialize("debug")
	s.Require().NoError(err)

	// Create a new Chi router and register a mock handler
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", s.mockMetricUpdateHandler)

	s.ts = httptest.NewServer(r)
	s.serverURL = s.ts.URL
}

func (s *MockAgentSuite) TearDownSuite() {
	s.ts.Close()
}

// mockMetricUpdateHandler is a mock HTTP handler that captures metrics sent to it.
func (s *MockAgentSuite) mockMetricUpdateHandler(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	s.mu.Lock()
	s.metrics = append(s.metrics, MockMetric{
		Type:  typ,
		Name:  name,
		Value: value,
	})
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (s *MockAgentSuite) TestAgentSendsMetrics() {
	cfg := &configs.AgentConfig{
		ServerAddress:  s.serverURL,
		ServerEndpoint: "/update/",
		LogLevel:       "debug",
		PollInterval:   1,
		ReportInterval: 1,
		NumWorkers:     1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := run(
		ctx,
		cfg,
		logger.Initialize,
		apps.NewAgentApp,
		runners.NewRunContext,
		runners.RunWorker,
	)

	// Accept context cancellation errors as normal termination
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		s.Require().NoError(err) // fail for any other error
	}

	time.Sleep(500 * time.Millisecond)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Require().NotEmpty(s.metrics)

	found := false
	for _, m := range s.metrics {
		if m.Name == "Alloc" && m.Type == "gauge" {
			found = true
			break
		}
	}
	s.Require().True(found)
}

func TestMockAgentSuite(t *testing.T) {
	suite.Run(t, new(MockAgentSuite))
}
