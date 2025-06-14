package types

import (
	"errors"
)

const (
	// Counter represents a metric type that counts occurrences.
	Counter = "counter"

	// Gauge represents a metric type that measures a value at a point in time.
	Gauge = "gauge"
)

var (
	// ErrInvalidMetricID indicates that a metric ID is invalid or empty.
	ErrInvalidMetricID = errors.New("invalid metric id")

	// ErrInvalidMetricType indicates that a metric type is invalid or unsupported.
	ErrInvalidMetricType = errors.New("invalid metric type")

	// ErrInvalidCounterValue indicates that a counter metric value is invalid.
	ErrInvalidCounterValue = errors.New("invalid counter value")

	// ErrInvalidGaugeValue indicates that a gauge metric value is invalid.
	ErrInvalidGaugeValue = errors.New("invalid gauge value")
)

// MetricID uniquely identifies a metric by its ID string and its type.
//
// Used as a key in maps and for metric identification in general.
type MetricID struct {
	ID    string `json:"id"`   // ID is the unique identifier for the metric.
	MType string `json:"type"` // MType is the metric type (e.g., counter, gauge).
}

// Metrics represents a metric with its ID, type, and optional values.
//
// Delta holds the value for counter metrics and Value for gauge metrics.
// Hash can be used for data integrity verification or similar purposes.
type Metrics struct {
	ID    string   `json:"id"`              // ID is the unique identifier for the metric.
	MType string   `json:"type"`            // MType is the metric type (counter or gauge).
	Delta *int64   `json:"delta,omitempty"` // Delta is the counter value, only used if MType is counter.
	Value *float64 `json:"value,omitempty"` // Value is the gauge value, only used if MType is gauge.
	Hash  string   `json:"hash,omitempty"`  // Hash is optional and can be used for validation or security.
}

// MetricsUpdatePathRequest представляет структуру запроса обновления метрики
// через URL-путь вида /update/{type}/{name}/{value}.
//
// Поля:
//   - Name:  имя метрики (например, "Alloc", "PollCount")
//   - MType: тип метрики ("gauge" или "counter")
//   - Value: значение метрики в строковом виде (float64 для gauge, int64 для counter)
type MetricsUpdatePathRequest struct {
	Name  string `json:"name"`  // Имя метрики
	MType string `json:"type"`  // Тип метрики: "gauge" или "counter"
	Value string `json:"value"` // Значение метрики (как строка)
}
