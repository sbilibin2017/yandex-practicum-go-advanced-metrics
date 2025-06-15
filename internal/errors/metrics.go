package errors

import "errors"

var (
	ErrInvalidMetricID     = errors.New("invalid metric id")
	ErrInvalidMetricType   = errors.New("invalid metric type")
	ErrInvalidCounterValue = errors.New("invalid counter value")
	ErrInvalidGaugeValue   = errors.New("invalid gauge value")
	ErrMetricNotFound      = errors.New("metric not found")
)
