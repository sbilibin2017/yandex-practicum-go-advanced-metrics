package workers

import (
	"context"
	"math/rand/v2"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricUpdater interface {
	Update(ctx context.Context, req types.MetricsUpdatePathRequest) error
}

func NewMetricAgentWorker(
	updater MetricUpdater,
	pollInterval int,
	reportInterval int,
	workerCount int,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return startMetricAgentWorker(ctx, updater, pollInterval, reportInterval, workerCount)
	}
}

func startMetricAgentWorker(
	ctx context.Context,
	updater MetricUpdater,
	pollInterval, reportInterval, workerCount int,
) error {
	collectors := []func() []types.MetricsUpdatePathRequest{
		collectRuntimeGaugeMetrics,
		collectRuntimeCounterMetrics,
	}
	metricsCh := pollMetrics(ctx, pollInterval, collectors...)
	errCh := reportMetrics(ctx, updater, reportInterval, workerCount, metricsCh)
	return waitForContextOrError(ctx, errCh)
}

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

func collectRuntimeCounterMetrics() []types.MetricsUpdatePathRequest {
	return []types.MetricsUpdatePathRequest{
		{MType: types.Counter, Name: "PollCount", Value: intToString(1)},
	}
}

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
				jobs <- metric
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

func waitForContextOrError(ctx context.Context, errCh <-chan error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
