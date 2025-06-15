package types

import (
	"html"
	"strconv"
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
	if metric == nil {
		return ""
	}

	switch metric.MType {
	case Counter:
		if metric.Delta == nil {
			return ""
		}
		return strconv.FormatInt(int64(*metric.Delta), 10)
	case Gauge:
		if metric.Value == nil {
			return ""
		}
		return strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	default:
		return ""
	}
}

func GetMetricsHTML(metrics []Metrics) string {
	htmlStr := "<html><head><title>Metrics List</title></head><body>"
	htmlStr += "<h1>Metrics</h1>"
	htmlStr += "<ul>"

	for _, metric := range metrics {
		name := html.EscapeString(metric.ID)
		value := html.EscapeString(GetMetricStringValue(&metric))
		htmlStr += "<li>" + name + ": " + value + "</li>"
	}

	htmlStr += "</ul></body></html>"
	return htmlStr
}

type MetricsUpdatePathRequest struct {
	Name  string `json:"name"`
	MType string `json:"type"`
	Value string `json:"value"`
}
