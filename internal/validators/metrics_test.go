package validators

import (
	"errors"
	"net/http"
	"testing"

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
		{"empty id", "", string(types.Counter), types.ErrInvalidMetricID},
		{"invalid type", "metric3", "invalid", types.ErrInvalidMetricType},
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
		{"empty id", "", string(types.Counter), "123", types.ErrInvalidMetricID},
		{"invalid type", "metric3", "invalid", "123", types.ErrInvalidMetricType},
		{"invalid counter value", "metric4", string(types.Counter), "abc", types.ErrInvalidCounterValue},
		{"invalid gauge value", "metric5", string(types.Gauge), "xyz", types.ErrInvalidGaugeValue},
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
			err:        ErrInvalidMetricID,
			wantStatus: http.StatusNotFound,
			wantMsg:    ErrInvalidMetricID.Error(),
		},
		{
			name:       "ErrInvalidMetricType returns 400",
			err:        ErrInvalidMetricType,
			wantStatus: http.StatusBadRequest,
			wantMsg:    ErrInvalidMetricType.Error(),
		},
		{
			name:       "ErrInvalidGaugeValue returns 400",
			err:        ErrInvalidGaugeValue,
			wantStatus: http.StatusBadRequest,
			wantMsg:    ErrInvalidGaugeValue.Error(),
		},
		{
			name:       "ErrInvalidCounterValue returns 400",
			err:        ErrInvalidCounterValue,
			wantStatus: http.StatusBadRequest,
			wantMsg:    ErrInvalidCounterValue.Error(),
		},
		{
			name:       "unknown error returns 500",
			err:        errors.New("some unknown error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HandleMetricsValidationError(tt.err)
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantStatus, got.Code)
			assert.Equal(t, tt.wantMsg, got.Message)
		})
	}
}
