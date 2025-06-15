package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricGetBodyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	float64Ptr := func(v float64) *float64 {
		return &v
	}

	mockSvc := NewMockMetricBodyGetter(ctrl)

	validMetricID := types.MetricID{ID: "testMetric", MType: types.Gauge} // Use MType field here
	expectedMetric := &types.Metrics{
		ID:    validMetricID.ID,
		MType: validMetricID.MType,
		Value: float64Ptr(1.23),
	}

	valFunc := func(id types.MetricID) error {
		if id.ID == "" {
			return errors.New("missing ID")
		}
		return nil
	}

	errValFunc := func(err error) *types.APIError {
		if err == nil {
			return nil
		}
		switch err.Error() {
		case "missing ID":
			return &types.APIError{Code: http.StatusBadRequest, Message: "missing ID"}
		case "service error":
			return &types.APIError{Code: http.StatusInternalServerError, Message: "service error"}
		default:
			return &types.APIError{Code: http.StatusInternalServerError, Message: "internal error"}
		}
	}

	handler := NewMetricGetBodyHandler(valFunc, errValFunc, mockSvc)

	t.Run("invalid content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/value/", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Content-Type")
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid JSON")
	})

	t.Run("validation error", func(t *testing.T) {
		payload := `{"id": "", "type": "gauge"}` // JSON key is "type"
		req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing ID")
	})

	t.Run("service error", func(t *testing.T) {
		payload := `{"id": "testMetric", "type": "gauge"}` // JSON key "type"
		req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		mockSvc.EXPECT().
			Get(gomock.Any(), gomock.Eq(validMetricID)). // Match MetricID with MType set
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})

	t.Run("success", func(t *testing.T) {
		payload, _ := json.Marshal(validMetricID) // will marshal with "type" key because of struct tag
		req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		mockSvc.EXPECT().
			Get(gomock.Any(), gomock.Eq(validMetricID)).
			Return(expectedMetric, nil)

		w := httptest.NewRecorder()
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp types.Metrics
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, *expectedMetric, resp)
	})
}

type failingResponseWriter struct{}

func (f *failingResponseWriter) Header() http.Header {
	return http.Header{}
}

func (f *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func (f *failingResponseWriter) WriteHeader(statusCode int) {}

func TestMetricGetBodyHandler_EncodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	float64Ptr := func(v float64) *float64 {
		return &v
	}

	mockSvc := NewMockMetricBodyGetter(ctrl)

	validMetricID := types.MetricID{ID: "testMetric", MType: "gauge"}
	expectedMetric := &types.Metrics{
		ID:    validMetricID.ID,
		MType: validMetricID.MType,
		Value: float64Ptr(1.23),
	}

	valFunc := func(id types.MetricID) error {
		return nil
	}

	errValFunc := func(err error) *types.APIError {
		return nil
	}

	handler := NewMetricGetBodyHandler(valFunc, errValFunc, mockSvc)

	payload, _ := json.Marshal(validMetricID)
	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	mockSvc.EXPECT().
		Get(gomock.Any(), gomock.Eq(validMetricID)).
		Return(expectedMetric, nil)

	w := &failingResponseWriter{}

	// This should call handleInternalServerError internally without panicking
	handler(w, req)
}
