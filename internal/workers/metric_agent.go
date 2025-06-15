package workers

import (
	"context"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricUpdater interface {
	Update(ctx context.Context, metrics types.Metrics) error
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
	collectors := []func() []types.Metrics{
		collectRuntimeGaugeMetrics,
		collectRuntimeCounterMetrics,
	}
	metricsCh := pollMetrics(ctx, pollInterval, collectors...)
	errCh := reportMetrics(ctx, updater, reportInterval, workerCount, metricsCh)
	return waitForContextOrError(ctx, errCh)
}

func collectRuntimeGaugeMetrics() []types.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	toPtrFloat := func(f float64) *float64 { return &f }

	return []types.Metrics{
		{MType: types.Gauge, ID: "Alloc", Value: toPtrFloat(float64(memStats.Alloc))},
		{MType: types.Gauge, ID: "BuckHashSys", Value: toPtrFloat(float64(memStats.BuckHashSys))},
		{MType: types.Gauge, ID: "Frees", Value: toPtrFloat(float64(memStats.Frees))},
		{MType: types.Gauge, ID: "GCCPUFraction", Value: toPtrFloat(memStats.GCCPUFraction)},
		{MType: types.Gauge, ID: "GCSys", Value: toPtrFloat(float64(memStats.GCSys))},
		{MType: types.Gauge, ID: "HeapAlloc", Value: toPtrFloat(float64(memStats.HeapAlloc))},
		{MType: types.Gauge, ID: "HeapIdle", Value: toPtrFloat(float64(memStats.HeapIdle))},
		{MType: types.Gauge, ID: "HeapInuse", Value: toPtrFloat(float64(memStats.HeapInuse))},
		{MType: types.Gauge, ID: "HeapObjects", Value: toPtrFloat(float64(memStats.HeapObjects))},
		{MType: types.Gauge, ID: "HeapReleased", Value: toPtrFloat(float64(memStats.HeapReleased))},
		{MType: types.Gauge, ID: "HeapSys", Value: toPtrFloat(float64(memStats.HeapSys))},
		{MType: types.Gauge, ID: "LastGC", Value: toPtrFloat(float64(memStats.LastGC))},
		{MType: types.Gauge, ID: "Lookups", Value: toPtrFloat(float64(memStats.Lookups))},
		{MType: types.Gauge, ID: "MCacheInuse", Value: toPtrFloat(float64(memStats.MCacheInuse))},
		{MType: types.Gauge, ID: "MCacheSys", Value: toPtrFloat(float64(memStats.MCacheSys))},
		{MType: types.Gauge, ID: "MSpanInuse", Value: toPtrFloat(float64(memStats.MSpanInuse))},
		{MType: types.Gauge, ID: "MSpanSys", Value: toPtrFloat(float64(memStats.MSpanSys))},
		{MType: types.Gauge, ID: "Mallocs", Value: toPtrFloat(float64(memStats.Mallocs))},
		{MType: types.Gauge, ID: "NextGC", Value: toPtrFloat(float64(memStats.NextGC))},
		{MType: types.Gauge, ID: "NumForcedGC", Value: toPtrFloat(float64(memStats.NumForcedGC))},
		{MType: types.Gauge, ID: "NumGC", Value: toPtrFloat(float64(memStats.NumGC))},
		{MType: types.Gauge, ID: "OtherSys", Value: toPtrFloat(float64(memStats.OtherSys))},
		{MType: types.Gauge, ID: "PauseTotalNs", Value: toPtrFloat(float64(memStats.PauseTotalNs))},
		{MType: types.Gauge, ID: "StackInuse", Value: toPtrFloat(float64(memStats.StackInuse))},
		{MType: types.Gauge, ID: "StackSys", Value: toPtrFloat(float64(memStats.StackSys))},
		{MType: types.Gauge, ID: "Sys", Value: toPtrFloat(float64(memStats.Sys))},
		{MType: types.Gauge, ID: "TotalAlloc", Value: toPtrFloat(float64(memStats.TotalAlloc))},
		{MType: types.Gauge, ID: "RandomValue", Value: toPtrFloat(rand.Float64() * 100)},
	}
}

func collectRuntimeCounterMetrics() []types.Metrics {
	v := int64(1)
	return []types.Metrics{
		{MType: types.Counter, ID: "PollCount", Delta: &v},
	}
}

func pollMetrics(
	ctx context.Context,
	pollInterval int,
	collectors ...func() []types.Metrics,
) <-chan types.Metrics {
	out := make(chan types.Metrics, 100)

	go func() {
		defer func() {
			log.Println("pollMetrics: stopping polling and closing channel")
			close(out)
		}()

		ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("pollMetrics: context canceled, exiting")
				return
			case <-ticker.C:
				log.Println("pollMetrics: polling metrics")
				for _, collect := range collectors {
					metrics := collect()
					log.Printf("pollMetrics: collected %d metrics", len(metrics))
					for _, metric := range metrics {
						log.Printf("pollMetrics: sending metric: ID=%s, Type=%s, Delta=%v, Value=%v",
							metric.ID, metric.MType, metric.Delta, metric.Value)
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
	in <-chan types.Metrics,
) <-chan error {
	errCh := make(chan error, 100)
	jobs := make(chan types.Metrics, 100)

	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		log.Printf("Worker %d: started", id)
		for metric := range jobs {
			log.Printf("Worker %d: updating metric ID=%s, Type=%s", id, metric.ID, metric.MType)
			if err := updater.Update(ctx, metric); err != nil {
				log.Printf("Worker %d: error updating metric ID=%s: %v", id, metric.ID, err)
				errCh <- err
			} else {
				log.Printf("Worker %d: successfully updated metric ID=%s", id, metric.ID)
			}
		}
		log.Printf("Worker %d: stopped", id)
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(i + 1)
	}

	go func() {
		defer func() {
			log.Println("Flusher: closing jobs channel")
			close(jobs)
		}()

		ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer ticker.Stop()

		var buffer []types.Metrics

		flush := func() {
			if len(buffer) == 0 {
				log.Println("Flusher: nothing to flush")
				return
			}
			log.Printf("Flusher: flushing %d metrics", len(buffer))
			for _, metric := range buffer {
				log.Printf("Flusher: sending metric ID=%s, Type=%s to jobs channel", metric.ID, metric.MType)
				jobs <- metric
			}
			buffer = buffer[:0]
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Flusher: context canceled, flushing remaining metrics and exiting")
				flush()
				return
			case metric, ok := <-in:
				if !ok {
					log.Println("Flusher: input channel closed, flushing remaining metrics and exiting")
					flush()
					return
				}
				log.Printf("Flusher: received metric ID=%s, Type=%s", metric.ID, metric.MType)
				buffer = append(buffer, metric)
			case <-ticker.C:
				log.Println("Flusher: ticker triggered flush")
				flush()
			}
		}
	}()

	go func() {
		wg.Wait()
		log.Println("All workers finished, closing error channel")
		close(errCh)
	}()

	return errCh
}

func waitForContextOrError(ctx context.Context, errCh <-chan error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			if err != nil {
				logger.Log.Error(err)
				return nil
			}
		}
	}
}
