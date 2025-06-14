package validators

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

var (
	// ErrInvalidMetricID indicates the metric ID is empty or invalid.
	ErrInvalidMetricID = errors.New("invalid metric id")

	// ErrInvalidMetricType indicates the metric type is not supported.
	ErrInvalidMetricType = errors.New("invalid metric type")

	// ErrInvalidCounterValue indicates the counter metric value is invalid.
	ErrInvalidCounterValue = errors.New("invalid counter value")

	// ErrInvalidGaugeValue indicates the gauge metric value is invalid.
	ErrInvalidGaugeValue = errors.New("invalid gauge value")
)

// ValidateMetricIDPath validates the metric ID and type for correctness.
// Returns an error if ID is empty or type is not one of the supported metric types.
func ValidateMetricIDPath(id string, mType string) error {
	if id == "" {
		logger.Log.Errorw("Validation failed: empty metric ID", "id", id)
		return ErrInvalidMetricID
	}

	if mType != types.Counter && mType != types.Gauge {
		logger.Log.Errorw("Validation failed: invalid metric type", "type", mType)
		return ErrInvalidMetricType
	}

	logger.Log.Debugw("Metric ID and type validated", "id", id, "type", mType)
	return nil
}

// ValidateMetricPath validates the metric ID, type, and string value.
// Returns an error if the ID or type is invalid or the value cannot be parsed
// to the expected type (int64 for Counter, float64 for Gauge).
func ValidateMetricPath(id string, mType string, value string) error {
	err := ValidateMetricIDPath(id, mType)
	if err != nil {
		return err
	}

	switch mType {
	case string(types.Counter):
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrInvalidCounterValue
		}
	case string(types.Gauge):
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrInvalidGaugeValue
		}
	}

	return nil
}

// HandleMetricsValidationError maps validation errors to appropriate API errors with HTTP status codes.
// Returns nil if no error is passed in.
func HandleMetricsValidationError(err error) *types.APIError {
	if err == nil {
		return nil
	}

	switch err {
	case ErrInvalidMetricID:
		return &types.APIError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	case ErrInvalidMetricType,
		ErrInvalidGaugeValue,
		ErrInvalidCounterValue:
		return &types.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	default:
		return &types.APIError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
}
