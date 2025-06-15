package types

import (
	"html"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricStringValue(t *testing.T) {
	counterDelta := int64(123)
	gaugeVal := float64(45.6789)

	tests := []struct {
		name     string
		metric   *Metrics
		expected string
	}{
		{
			name: "Counter type with Delta set",
			metric: &Metrics{
				MType: Counter,
				Delta: &counterDelta,
			},
			expected: strconv.FormatInt(counterDelta, 10),
		},
		{
			name: "Gauge type with Value set",
			metric: &Metrics{
				MType: Gauge,
				Value: &gaugeVal,
			},
			expected: strconv.FormatFloat(gaugeVal, 'f', -1, 64),
		},
		{
			name:     "Nil metric returns empty string",
			metric:   nil,
			expected: "",
		},
		{
			name: "Counter type with nil Delta returns empty string",
			metric: &Metrics{
				MType: Counter,
				Delta: nil,
			},
			expected: "",
		},
		{
			name: "Gauge type with nil Value returns empty string",
			metric: &Metrics{
				MType: Gauge,
				Value: nil,
			},
			expected: "",
		},
		{
			name: "Unknown type returns empty string",
			metric: &Metrics{
				MType: "unknown",
				Value: &gaugeVal,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMetricStringValue(tt.metric)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMetricsHTML(t *testing.T) {
	counterVal := float64(123)
	gaugeVal := float64(45.6789)

	metrics := []Metrics{
		{
			ID:    "metric1",
			MType: Counter,
			Value: &counterVal,
		},
		{
			ID:    "metric2",
			MType: Gauge,
			Value: &gaugeVal,
		},
		{
			ID:    "metric<3",
			MType: Gauge,
			Value: &gaugeVal,
		},
	}

	result := GetMetricsHTML(metrics)

	// Basic checks
	assert.Contains(t, result, "<html>")
	assert.Contains(t, result, "<h1>Metrics</h1>")
	assert.Contains(t, result, "<ul>")
	assert.Contains(t, result, "</ul>")
	assert.Contains(t, result, "</html>")

	// Check escaped metric ID
	assert.Contains(t, result, "metric1: "+GetMetricStringValue(&metrics[0]))
	assert.Contains(t, result, "metric2: "+GetMetricStringValue(&metrics[1]))

	// Check HTML escaping (metric ID "metric<3" should be escaped)
	escapedID := html.EscapeString(metrics[2].ID)
	expectedEntry := "<li>" + escapedID + ": " + GetMetricStringValue(&metrics[2]) + "</li>"
	assert.Contains(t, result, expectedEntry)
}

func TestGetMetricsHTMLEmpty(t *testing.T) {
	result := GetMetricsHTML(nil)
	assert.Contains(t, result, "<ul></ul>")
	assert.Contains(t, result, "<h1>Metrics</h1>")
}
