package workers

import (
	"context"
	"math/rand/v2"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricUpdater defines an interface for sending metric updates.
//
// Implementations should provide logic for sending metrics, e.g., via HTTP.
type MetricUpdater interface {
	// Update sends a metric update request.
	//
	// Parameters:
	//   - ctx: context for cancellation and timeout.
	//   - req: the metric update request to send.
	//
	// Returns:
	//   - error: non-nil if sending failed.
	Update(ctx context.Context, req types.MetricsUpdatePathRequest) error
}

// NewMetricAgentWorker returns a worker function that runs metric polling
// and reporting at specified intervals.
//
// Parameters:
//   - updater: implementation of MetricUpdater to send metrics.
//   - pollInterval: seconds between metric collection.
//   - reportInterval: seconds between sending collected metrics.
//   - workerCount: number of concurrent workers to report metrics.
//
// Returns:
//   - func(ctx context.Context): worker function to be run with a context.
func NewMetricAgentWorker(
	updater MetricUpdater,
	pollInterval int,
	reportInterval int,
	workerCount int,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		startMetricAgentWorker(ctx, updater, pollInterval, reportInterval, workerCount)
	}
}

// startMetricAgentWorker runs the main worker loop that polls and reports metrics.
//
// Parameters:
//   - ctx: context for lifecycle control and cancellation.
//   - updater: MetricUpdater to send metric updates.
//   - pollInterval: interval in seconds between polling metrics.
//   - reportInterval: interval in seconds between reporting metrics.
//   - workerCount: number of concurrent workers processing reports.
func startMetricAgentWorker(ctx context.Context, updater MetricUpdater, pollInterval, reportInterval, workerCount int) {
	collectors := []func() []types.MetricsUpdatePathRequest{
		collectRuntimeGaugeMetrics,
		collectRuntimeCounterMetrics,
	}

	metricsCh := pollMetrics(ctx, pollInterval, collectors...)
	errCh := reportMetrics(ctx, updater, reportInterval, workerCount, metricsCh)
	logErrors(ctx, errCh)
}

// collectRuntimeGaugeMetrics gathers runtime memory statistics
// and returns them as a slice of gauge-type metrics.
//
// Returns:
//   - []types.MetricsUpdatePathRequest: list of gauge metrics with current runtime values.
func collectRuntimeGaugeMetrics() []types.MetricsUpdatePathRequest {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return []types.MetricsUpdatePathRequest{
		{MType: types.Gauge, Name: "Alloc", Value: floatToString(float64(memStats.Alloc))},
		{MType: types.Gauge, Name: "BuckHashSys", Value: floatToString(float64(memStats.BuckHashSys))},
		{MType: types.Gauge, Name: "Frees", Value: floatToString(float64(memStats.Frees))},
		{MType: types.Gauge, Name: "GCCPUFraction", Value: floatToString(memStats.GCCPUFraction)},
		{MType: types.Gauge, Name: "GCSys", Value: floatToString(float64(memStats.GCSys))},
		{MType: types.Gauge, Name: "HeapAlloc", Value: floatToString(float64(memStats.HeapAlloc))},
		{MType: types.Gauge, Name: "HeapIdle", Value: floatToString(float64(memStats.HeapIdle))},
		{MType: types.Gauge, Name: "HeapInuse", Value: floatToString(float64(memStats.HeapInuse))},
		{MType: types.Gauge, Name: "HeapObjects", Value: floatToString(float64(memStats.HeapObjects))},
		{MType: types.Gauge, Name: "HeapReleased", Value: floatToString(float64(memStats.HeapReleased))},
		{MType: types.Gauge, Name: "HeapSys", Value: floatToString(float64(memStats.HeapSys))},
		{MType: types.Gauge, Name: "LastGC", Value: floatToString(float64(memStats.LastGC))},
		{MType: types.Gauge, Name: "Lookups", Value: floatToString(float64(memStats.Lookups))},
		{MType: types.Gauge, Name: "MCacheInuse", Value: floatToString(float64(memStats.MCacheInuse))},
		{MType: types.Gauge, Name: "MCacheSys", Value: floatToString(float64(memStats.MCacheSys))},
		{MType: types.Gauge, Name: "MSpanInuse", Value: floatToString(float64(memStats.MSpanInuse))},
		{MType: types.Gauge, Name: "MSpanSys", Value: floatToString(float64(memStats.MSpanSys))},
		{MType: types.Gauge, Name: "Mallocs", Value: floatToString(float64(memStats.Mallocs))},
		{MType: types.Gauge, Name: "NextGC", Value: floatToString(float64(memStats.NextGC))},
		{MType: types.Gauge, Name: "NumForcedGC", Value: floatToString(float64(memStats.NumForcedGC))},
		{MType: types.Gauge, Name: "NumGC", Value: floatToString(float64(memStats.NumGC))},
		{MType: types.Gauge, Name: "OtherSys", Value: floatToString(float64(memStats.OtherSys))},
		{MType: types.Gauge, Name: "PauseTotalNs", Value: floatToString(float64(memStats.PauseTotalNs))},
		{MType: types.Gauge, Name: "StackInuse", Value: floatToString(float64(memStats.StackInuse))},
		{MType: types.Gauge, Name: "StackSys", Value: floatToString(float64(memStats.StackSys))},
		{MType: types.Gauge, Name: "Sys", Value: floatToString(float64(memStats.Sys))},
		{MType: types.Gauge, Name: "TotalAlloc", Value: floatToString(float64(memStats.TotalAlloc))},
		{MType: types.Gauge, Name: "RandomValue", Value: floatToString(rand.Float64() * 100)},
	}
}

// collectRuntimeCounterMetrics returns runtime counter metrics.
//
// Returns:
//   - []types.MetricsUpdatePathRequest: list containing counter metrics.
func collectRuntimeCounterMetrics() []types.MetricsUpdatePathRequest {
	return []types.MetricsUpdatePathRequest{
		{MType: types.Counter, Name: "PollCount", Value: intToString(1)},
	}
}

// pollMetrics periodically collects metrics using the provided collector functions.
//
// Parameters:
//   - ctx: context for cancellation.
//   - pollInterval: interval in seconds between polls.
//   - collectors: variadic list of functions that return metrics.
//
// Returns:
//   - <-chan types.MetricsUpdatePathRequest: channel streaming collected metrics.
func pollMetrics(
	ctx context.Context,
	pollInterval int,
	collectors ...func() []types.MetricsUpdatePathRequest,
) <-chan types.MetricsUpdatePathRequest {
	out := make(chan types.MetricsUpdatePathRequest, 100)

	go func() {
		defer close(out)
		ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, collect := range collectors {
					metrics := collect()
					for _, metric := range metrics {
						out <- metric
					}
				}
			}
		}
	}()

	return out
}

// reportMetrics buffers and sends metrics through a worker pool.
//
// Parameters:
//   - ctx: context for cancellation.
//   - updater: interface used to send metrics.
//   - reportInterval: flush interval in seconds.
//   - workerCount: number of concurrent workers.
//   - in: channel from which metrics are received.
//
// Returns:
//   - <-chan error: channel streaming any update errors.
func reportMetrics(
	ctx context.Context,
	updater MetricUpdater,
	reportInterval int,
	workerCount int,
	in <-chan types.MetricsUpdatePathRequest,
) <-chan error {
	errCh := make(chan error, 100)
	jobs := make(chan types.MetricsUpdatePathRequest, 100)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for metric := range jobs {
			if err := updater.Update(ctx, metric); err != nil {
				errCh <- err
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	go func() {
		defer close(jobs)
		ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer ticker.Stop()

		var buffer []types.MetricsUpdatePathRequest

		flush := func() {
			for _, metric := range buffer {
				select {
				case jobs <- metric:
				case <-ctx.Done():
					return
				}
			}
			buffer = buffer[:0]
		}

		for {
			select {
			case <-ctx.Done():
				flush()
				return
			case metric, ok := <-in:
				if !ok {
					flush()
					return
				}
				buffer = append(buffer, metric)
			case <-ticker.C:
				flush()
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}

// logErrors listens to the error channel and logs errors until the context is canceled.
//
// Parameters:
//   - ctx: context for cancellation.
//   - errCh: channel from which errors are received.
func logErrors(ctx context.Context, errCh <-chan error) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stopped logging errors due to context cancellation")
			return
		case err, ok := <-errCh:
			if !ok {
				logger.Log.Info("Error channel closed, stopping error logging")
				return
			}
			if err != nil {
				logger.Log.Errorf("Metric update error: %v", err)
			}
		}
	}
}

// intToString converts an int64 to its string representation.
//
// Parameters:
//   - i: the integer to convert.
//
// Returns:
//   - string representation of i.
func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// floatToString converts a float64 to its string representation.
//
// Parameters:
//   - f: the float to convert.
//
// Returns:
//   - string representation of f.
func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
