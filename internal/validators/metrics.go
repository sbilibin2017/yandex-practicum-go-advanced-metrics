package validators

import (
	"net/http"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func ValidateMetricIDPath(id string, mType string) error {
	logger.Log.Debugf("ValidateMetricIDPath called with id=%q, mType=%q", id, mType)

	if id == "" {
		logger.Log.Debugf("Validation error: invalid metric ID")
		return errors.ErrInvalidMetricID
	}

	if mType != types.Counter && mType != types.Gauge {
		logger.Log.Debugf("Validation error: invalid metric type %q", mType)
		return errors.ErrInvalidMetricType
	}

	return nil
}

func ValidateMetricPath(id string, mType string, value string) error {
	logger.Log.Debugf("ValidateMetricPath called with id=%q, mType=%q, value=%q", id, mType, value)

	err := ValidateMetricIDPath(id, mType)
	if err != nil {
		logger.Log.Debugf("ValidateMetricPath validation failed: %v", err)
		return err
	}

	switch mType {
	case string(types.Counter):
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Log.Debugf("Validation error: invalid counter value %q", value)
			return errors.ErrInvalidCounterValue
		}
	case string(types.Gauge):
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Log.Debugf("Validation error: invalid gauge value %q", value)
			return errors.ErrInvalidGaugeValue
		}
	}

	return nil
}

func ValidateMetricBody(metric types.Metrics) error {
	logger.Log.Debugf("ValidateMetricBody called with metric ID=%q, Type=%q", metric.ID, metric.MType)

	err := ValidateMetricIDPath(metric.ID, metric.MType)
	if err != nil {
		logger.Log.Debugf("ValidateMetricBody validation failed: %v", err)
		return err
	}

	switch metric.MType {
	case string(types.Counter):
		if metric.Delta == nil {
			logger.Log.Debugf("Validation error: counter metric Delta is nil")
			return errors.ErrInvalidCounterValue
		}
	case string(types.Gauge):
		if metric.Value == nil {
			logger.Log.Debugf("Validation error: gauge metric Value is nil")
			return errors.ErrInvalidGaugeValue
		}
	}

	return nil
}

func ValidateMetricIDBody(id types.MetricID) error {
	logger.Log.Debugf("ValidateMetricIDBody called with ID=%q, Type=%q", id.ID, id.MType)
	err := ValidateMetricIDPath(id.ID, id.MType)
	if err != nil {
		logger.Log.Debugf("ValidateMetricIDBody validation failed: %v", err)
		return err
	}
	return nil
}

func HandleMetricsValidationError(err error) *types.APIError {
	if err == nil {
		return nil
	}

	logger.Log.Debugf("HandleMetricsValidationError called with error: %v", err)

	switch err {
	case errors.ErrInvalidMetricID,
		errors.ErrMetricNotFound:
		return &types.APIError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	case errors.ErrInvalidMetricType,
		errors.ErrInvalidGaugeValue,
		errors.ErrInvalidCounterValue:
		return &types.APIError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	default:
		return &types.APIError{
			Code:    http.StatusInternalServerError,
			Message: errors.ErrInternalServerError.Error(),
		}
	}
}
