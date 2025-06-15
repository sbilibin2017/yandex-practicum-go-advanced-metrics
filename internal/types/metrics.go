package types

import (
	"html"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MetricID struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func GetMetricStringValue(metric *Metrics) string {
	logger.Log.Debugf("GetMetricStringValue called with metric: %+v", metric)

	if metric == nil {
		logger.Log.Debug("Metric is nil, returning empty string")
		return ""
	}

	switch metric.MType {
	case Counter:
		if metric.Delta == nil {
			logger.Log.Debug("Metric type is Counter but Delta is nil, returning empty string")
			return ""
		}
		val := strconv.FormatInt(int64(*metric.Delta), 10)
		logger.Log.Debugf("Metric type Counter, value: %s", val)
		return val
	case Gauge:
		if metric.Value == nil {
			logger.Log.Debug("Metric type is Gauge but Value is nil, returning empty string")
			return ""
		}
		val := strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		logger.Log.Debugf("Metric type Gauge, value: %s", val)
		return val
	default:
		logger.Log.Debugf("Unknown metric type %q, returning empty string", metric.MType)
		return ""
	}
}

func GetMetricsHTML(metrics []Metrics) string {
	logger.Log.Debugf("GetMetricsHTML called with %d metrics", len(metrics))

	htmlStr := "<html><head><title>Metrics List</title></head><body>"
	htmlStr += "<h1>Metrics</h1>"
	htmlStr += "<ul>"

	for _, metric := range metrics {
		name := html.EscapeString(metric.ID)
		value := html.EscapeString(GetMetricStringValue(&metric))
		logger.Log.Debugf("Adding metric to HTML: %s = %s", name, value)
		htmlStr += "<li>" + name + ": " + value + "</li>"
	}

	htmlStr += "</ul></body></html>"

	logger.Log.Debug("Generated HTML metrics list")
	return htmlStr
}
