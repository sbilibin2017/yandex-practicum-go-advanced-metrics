package validators

import (
	"net/http"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func ValidateMetricIDPath(id string, mType string) error {
	if id == "" {
		return errors.ErrInvalidMetricID
	}

	if mType != types.Counter && mType != types.Gauge {
		return errors.ErrInvalidMetricType
	}

	return nil
}

func ValidateMetricPath(id string, mType string, value string) error {
	err := ValidateMetricIDPath(id, mType)
	if err != nil {
		return err
	}

	switch mType {
	case string(types.Counter):
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errors.ErrInvalidCounterValue
		}
	case string(types.Gauge):
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.ErrInvalidGaugeValue
		}
	}

	return nil
}

func HandleMetricsValidationError(err error) *types.APIError {
	if err == nil {
		return nil
	}

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
