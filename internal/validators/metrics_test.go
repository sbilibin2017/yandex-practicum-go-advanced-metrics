package validators

import (
	"errors"
	"net/http"
	"testing"

	internalErrors "github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateMetricIDPath(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mType   string
		wantErr error
	}{
		{"valid counter", "metric1", string(types.Counter), nil},
		{"valid gauge", "metric2", string(types.Gauge), nil},
		{"empty id", "", string(types.Counter), internalErrors.ErrInvalidMetricID},
		{"invalid type", "metric3", "invalid", internalErrors.ErrInvalidMetricType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricIDPath(tt.id, tt.mType)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestValidateMetricPath(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mType   string
		value   string
		wantErr error
	}{
		{"valid counter", "metric1", string(types.Counter), "123", nil},
		{"valid gauge", "metric2", string(types.Gauge), "123.456", nil},
		{"empty id", "", string(types.Counter), "123", internalErrors.ErrInvalidMetricID},
		{"invalid type", "metric3", "invalid", "123", internalErrors.ErrInvalidMetricType},
		{"invalid counter value", "metric4", string(types.Counter), "abc", internalErrors.ErrInvalidCounterValue},
		{"invalid gauge value", "metric5", string(types.Gauge), "xyz", internalErrors.ErrInvalidGaugeValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricPath(tt.id, tt.mType, tt.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHandleMetricsValidationError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
		wantNil    bool
	}{
		{
			name:    "nil error returns nil",
			err:     nil,
			wantNil: true,
		},
		{
			name:       "ErrInvalidMetricID returns 404",
			err:        internalErrors.ErrInvalidMetricID,
			wantStatus: http.StatusNotFound,
			wantMsg:    internalErrors.ErrInvalidMetricID.Error(),
		},
		{
			name:       "ErrInvalidMetricType returns 400",
			err:        internalErrors.ErrInvalidMetricType,
			wantStatus: http.StatusBadRequest,
			wantMsg:    internalErrors.ErrInvalidMetricType.Error(),
		},
		{
			name:       "ErrInvalidGaugeValue returns 400",
			err:        internalErrors.ErrInvalidGaugeValue,
			wantStatus: http.StatusBadRequest,
			wantMsg:    internalErrors.ErrInvalidGaugeValue.Error(),
		},
		{
			name:       "ErrInvalidCounterValue returns 400",
			err:        internalErrors.ErrInvalidCounterValue,
			wantStatus: http.StatusBadRequest,
			wantMsg:    internalErrors.ErrInvalidCounterValue.Error(),
		},
		{
			name:       "unknown error returns 500",
			err:        errors.New("some unknown error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    internalErrors.ErrInternalServerError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HandleMetricsValidationError(tt.err)
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantStatus, got.Code)
				assert.Equal(t, tt.wantMsg, got.Message)
			}
		})
	}
}

func TestValidateMetricBody(t *testing.T) {
	counterDelta := int64(100)
	gaugeValue := 12.34

	tests := []struct {
		name    string
		metric  types.Metrics
		wantErr error
	}{
		{
			name: "valid counter metric",
			metric: types.Metrics{
				ID:    "metric1",
				MType: string(types.Counter),
				Delta: &counterDelta,
			},
			wantErr: nil,
		},
		{
			name: "invalid counter metric missing delta",
			metric: types.Metrics{
				ID:    "metric2",
				MType: string(types.Counter),
				Delta: nil,
			},
			wantErr: internalErrors.ErrInvalidCounterValue,
		},
		{
			name: "valid gauge metric",
			metric: types.Metrics{
				ID:    "metric3",
				MType: string(types.Gauge),
				Value: &gaugeValue,
			},
			wantErr: nil,
		},
		{
			name: "invalid gauge metric missing value",
			metric: types.Metrics{
				ID:    "metric4",
				MType: string(types.Gauge),
				Value: nil,
			},
			wantErr: internalErrors.ErrInvalidGaugeValue,
		},
		{
			name: "invalid metric id (empty)",
			metric: types.Metrics{
				ID:    "",
				MType: string(types.Counter),
				Delta: &counterDelta,
			},
			wantErr: internalErrors.ErrInvalidMetricID,
		},
		{
			name: "invalid metric type",
			metric: types.Metrics{
				ID:    "metric6",
				MType: "invalid",
				Delta: &counterDelta,
			},
			wantErr: internalErrors.ErrInvalidMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricBody(tt.metric)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestValidateMetricIDBody(t *testing.T) {
	tests := []struct {
		name    string
		id      types.MetricID
		wantErr error
	}{
		{
			name:    "valid counter",
			id:      types.MetricID{ID: "metric1", MType: string(types.Counter)},
			wantErr: nil,
		},
		{
			name:    "valid gauge",
			id:      types.MetricID{ID: "metric2", MType: string(types.Gauge)},
			wantErr: nil,
		},
		{
			name:    "empty id",
			id:      types.MetricID{ID: "", MType: string(types.Counter)},
			wantErr: internalErrors.ErrInvalidMetricID,
		},
		{
			name:    "invalid type",
			id:      types.MetricID{ID: "metric3", MType: "invalid"},
			wantErr: internalErrors.ErrInvalidMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricIDBody(tt.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
