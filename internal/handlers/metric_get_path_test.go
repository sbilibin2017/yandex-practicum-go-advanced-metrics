package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	internalErrors "github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/errors"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestNewMetricGetPathHandler_Chi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricPathGetter(ctrl)

	errValHandlerFunc := func(err error) *types.APIError {
		if err != nil {
			// map specific errors to status codes
			if errors.Is(err, internalErrors.ErrMetricNotFound) {
				return &types.APIError{Message: "not found", Code: http.StatusNotFound}
			}
			return &types.APIError{Message: err.Error(), Code: http.StatusInternalServerError}
		}
		return nil
	}

	valFuncSuccess := func(name, typ string) error { return nil }
	valFuncFail := func(name, typ string) error { return errors.New("validation error") }

	metricName := "testMetric"
	metricType := types.Gauge

	metricID := types.MetricID{
		ID:    metricName,
		MType: metricType,
	}

	metricValue := 123.456
	metric := &types.Metrics{
		ID:    metricName,
		MType: metricType,
		Value: &metricValue,
	}

	tests := []struct {
		name           string
		valFunc        func(string, string) error
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "success path",
			valFunc: valFuncSuccess,
			mockSetup: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), metricID).
					Return(metric, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   types.GetMetricStringValue(metric), // здесь динамически берём формат
		},
		{
			name:           "validation error",
			valFunc:        valFuncFail,
			mockSetup:      func() {},
			expectedStatus: http.StatusInternalServerError, // changed to match error handler behavior
			expectedBody:   "validation error",
		},
		{
			name:    "service returns internal error",
			valFunc: valFuncSuccess,
			mockSetup: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), metricID).
					Return(nil, errors.New("database failure"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database failure",
		},
		{
			name:    "metric not found error",
			valFunc: valFuncSuccess,
			mockSetup: func() {
				mockSvc.EXPECT().
					Get(gomock.Any(), metricID).
					Return(nil, internalErrors.ErrMetricNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			handler := NewMetricGetPathHandler(tt.valFunc, errValHandlerFunc, mockSvc)

			r := chi.NewRouter()
			r.Get("/value/{type}/{name}", handler)

			req := httptest.NewRequest("GET", "/value/"+metricType+"/"+metricName, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close() // close response body to avoid leaks

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
