package workers

// // Test collectRuntimeGaugeMetrics returns expected metrics with correct types and values.
// func TestCollectRuntimeGaugeMetrics(t *testing.T) {
// 	metrics := collectRuntimeGaugeMetrics()

// 	assert.NotEmpty(t, metrics)
// 	for _, m := range metrics {
// 		assert.Equal(t, types.Gauge, m.MType)
// 		assert.NotEmpty(t, m.ID)
// 		assert.NotNil(t, m.Value)
// 		assert.Nil(t, m.Delta)
// 	}
// }

// // Test collectRuntimeCounterMetrics returns expected metric.
// func TestCollectRuntimeCounterMetrics(t *testing.T) {
// 	metrics := collectRuntimeCounterMetrics()

// 	assert.Len(t, metrics, 1)
// 	m := metrics[0]
// 	assert.Equal(t, types.Counter, m.MType)
// 	assert.Equal(t, "PollCount", m.ID)
// 	assert.NotNil(t, m.Delta)
// 	assert.Nil(t, m.Value)
// 	assert.Equal(t, int64(1), *m.Delta)
// }

// // Test pollMetrics emits metrics periodically until context canceled.
// func TestPollMetrics(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	pollInterval := 1
// 	collectors := []func() []types.Metrics{
// 		func() []types.Metrics {
// 			val := 42.0
// 			return []types.Metrics{
// 				{MType: types.Gauge, ID: "test_metric", Value: &val},
// 			}
// 		},
// 	}

// 	ch := pollMetrics(ctx, pollInterval, collectors...)

// 	// Read first metric, assert correctness
// 	select {
// 	case metric := <-ch:
// 		assert.Equal(t, "test_metric", metric.ID)
// 		assert.Equal(t, types.Gauge, metric.MType)
// 		assert.NotNil(t, metric.Value)
// 	case <-time.After(2 * time.Second):
// 		t.Fatal("Timeout waiting for metric")
// 	}

// 	// Cancel context and ensure channel closes
// 	cancel()
// 	time.Sleep(100 * time.Millisecond)
// 	_, ok := <-ch
// 	assert.False(t, ok)
// }

// func TestReportMetrics(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUpdater := NewMockMetricUpdater(ctrl)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	inCh := make(chan types.Metrics, 3)
// 	metrics := []types.Metrics{
// 		{MType: types.Counter, ID: "m1", Delta: int64Ptr(1)},
// 		{MType: types.Counter, ID: "m2", Delta: int64Ptr(2)},
// 		{MType: types.Counter, ID: "m3", Delta: int64Ptr(3)},
// 	}
// 	for _, m := range metrics {
// 		inCh <- m
// 	}
// 	close(inCh)

// 	for _, m := range metrics {
// 		mockUpdater.EXPECT().Update(gomock.Any(), m).Return(nil).Times(1)
// 	}

// 	errCh := reportMetrics(ctx, mockUpdater, 1, 2, inCh)

// 	// Wait a bit more than the reportInterval to let metrics flush
// 	time.Sleep(1100 * time.Millisecond)

// 	cancel()

// 	for err := range errCh {
// 		assert.NoError(t, err)
// 	}
// }

// // Test reportMetrics propagates errors from updater.
// func TestReportMetricsWithError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUpdater := NewMockMetricUpdater(ctrl)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	inCh := make(chan types.Metrics, 1)
// 	metric := types.Metrics{MType: types.Counter, ID: "errMetric", Delta: int64Ptr(1)}
// 	inCh <- metric
// 	close(inCh)

// 	expectedErr := errors.New("update failed")

// 	mockUpdater.EXPECT().Update(gomock.Any(), metric).Return(expectedErr)

// 	errCh := reportMetrics(ctx, mockUpdater, 1, 1, inCh)

// 	select {
// 	case err := <-errCh:
// 		assert.EqualError(t, err, expectedErr.Error())
// 	case <-time.After(2 * time.Second):
// 		t.Fatal("Timeout waiting for error")
// 	}
// }

// // Test NewMetricAgentWorker runs and respects context cancellation.
// func TestNewMetricAgentWorker(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUpdater := NewMockMetricUpdater(ctrl)

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Expect updater.Update to be called at least once
// 	mockUpdater.EXPECT().Update(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

// 	worker := NewMetricAgentWorker(mockUpdater, 1, 1, 2)

// 	doneCh := make(chan struct{})
// 	go func() {
// 		worker(ctx)
// 		close(doneCh)
// 	}()

// 	select {
// 	case <-doneCh:
// 		// worker exited normally
// 	case <-time.After(4 * time.Second):
// 		t.Fatal("Worker did not stop after context cancellation")
// 	}
// }

// func TestWaitForContextOrError(t *testing.T) {
// 	t.Run("returns context error when context is done", func(t *testing.T) {
// 		ctx, cancel := context.WithCancel(context.Background())
// 		cancel() // cancel immediately

// 		errCh := make(chan error)

// 		err := waitForContextOrError(ctx, errCh)
// 		require.ErrorIs(t, err, context.Canceled)
// 	})

// 	t.Run("returns error from channel", func(t *testing.T) {
// 		ctx := context.Background()
// 		errCh := make(chan error, 1)
// 		expectedErr := errors.New("something went wrong")

// 		errCh <- expectedErr
// 		close(errCh)

// 		err := waitForContextOrError(ctx, errCh)
// 		require.EqualError(t, err, expectedErr.Error())
// 	})

// 	t.Run("returns nil when channel is closed with no errors", func(t *testing.T) {
// 		ctx := context.Background()
// 		errCh := make(chan error)
// 		close(errCh)

// 		err := waitForContextOrError(ctx, errCh)
// 		require.NoError(t, err)
// 	})

// 	t.Run("blocks until context done if no error received", func(t *testing.T) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
// 		defer cancel()

// 		errCh := make(chan error)

// 		start := time.Now()
// 		err := waitForContextOrError(ctx, errCh)
// 		duration := time.Since(start)

// 		require.ErrorIs(t, err, context.DeadlineExceeded)
// 		require.GreaterOrEqual(t, duration.Milliseconds(), int64(50))
// 	})
// }

// // helper to get pointer to int64 literal
// func int64Ptr(i int64) *int64 {
// 	return &i
// }
