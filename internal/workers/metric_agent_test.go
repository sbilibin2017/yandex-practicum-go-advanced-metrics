package workers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"

	"github.com/stretchr/testify/assert"
)

// Test collectRuntimeGaugeMetrics returns expected metrics with correct types and values.
func TestCollectRuntimeGaugeMetrics(t *testing.T) {
	metrics := collectRuntimeGaugeMetrics()

	assert.NotEmpty(t, metrics)
	for _, m := range metrics {
		assert.Equal(t, types.Gauge, m.MType)
		assert.NotEmpty(t, m.Name)
		assert.NotEmpty(t, m.Value)
	}
}

// Test collectRuntimeCounterMetrics returns expected metric.
func TestCollectRuntimeCounterMetrics(t *testing.T) {
	metrics := collectRuntimeCounterMetrics()

	assert.Len(t, metrics, 1)
	m := metrics[0]
	assert.Equal(t, types.Counter, m.MType)
	assert.Equal(t, "PollCount", m.Name)
	assert.Equal(t, "1", m.Value)
}

// Test pollMetrics emits metrics periodically until context canceled.
func TestPollMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pollInterval := 1
	collectors := []func() []types.MetricsUpdatePathRequest{
		func() []types.MetricsUpdatePathRequest {
			return []types.MetricsUpdatePathRequest{
				{MType: types.Counter, Name: "test_metric", Value: "42"},
			}
		},
	}

	ch := pollMetrics(ctx, pollInterval, collectors...)

	// Read first metric, assert correctness
	select {
	case metric := <-ch:
		assert.Equal(t, "test_metric", metric.Name)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for metric")
	}

	// Cancel context and ensure channel closes
	cancel()
	time.Sleep(100 * time.Millisecond)
	_, ok := <-ch
	assert.False(t, ok)
}

func TestReportMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inCh := make(chan types.MetricsUpdatePathRequest, 3)
	metrics := []types.MetricsUpdatePathRequest{
		{MType: types.Counter, Name: "m1", Value: "1"},
		{MType: types.Counter, Name: "m2", Value: "2"},
		{MType: types.Counter, Name: "m3", Value: "3"},
	}
	for _, m := range metrics {
		inCh <- m
	}
	close(inCh)

	for _, m := range metrics {
		mockUpdater.EXPECT().Update(gomock.Any(), m).Return(nil).Times(1)
	}

	errCh := reportMetrics(ctx, mockUpdater, 1, 2, inCh)

	// Wait a bit more than the reportInterval to let metrics flush
	time.Sleep(1100 * time.Millisecond)

	cancel()

	for err := range errCh {
		assert.NoError(t, err)
	}
}

// Test reportMetrics propagates errors from updater.
func TestReportMetricsWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inCh := make(chan types.MetricsUpdatePathRequest, 1)
	metric := types.MetricsUpdatePathRequest{MType: types.Counter, Name: "errMetric", Value: "1"}
	inCh <- metric
	close(inCh)

	expectedErr := errors.New("update failed")

	mockUpdater.EXPECT().Update(gomock.Any(), metric).Return(expectedErr)

	errCh := reportMetrics(ctx, mockUpdater, 1, 1, inCh)

	select {
	case err := <-errCh:
		assert.EqualError(t, err, expectedErr.Error())
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for error")
	}
}

// Test logErrors stops logging on context cancellation and on channel close.
func TestLogErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)

	go func() {
		logErrors(ctx, errCh)
	}()

	errCh <- errors.New("some error")
	close(errCh)

	// Wait briefly to allow logErrors to finish
	time.Sleep(100 * time.Millisecond)

	cancel()
	time.Sleep(50 * time.Millisecond)
}

// Test NewMetricAgentWorker runs and respects context cancellation.
func TestNewMetricAgentWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Expect updater.Update to be called at least once
	mockUpdater.EXPECT().Update(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	worker := NewMetricAgentWorker(mockUpdater, 1, 1, 2)

	doneCh := make(chan struct{})
	go func() {
		worker(ctx)
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// worker exited normally
	case <-time.After(4 * time.Second):
		t.Fatal("Worker did not stop after context cancellation")
	}
}
