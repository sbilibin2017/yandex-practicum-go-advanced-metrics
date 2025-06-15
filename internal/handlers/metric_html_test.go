package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestNewMetricListHTMLHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricHTMLLister(ctrl)

	errHandlerFunc := func(err error) *types.APIError {
		if err != nil {
			return &types.APIError{Message: err.Error(), Code: http.StatusInternalServerError}
		}
		return nil
	}

	metricValue := 42.0
	metrics := []types.Metrics{
		{
			ID:    "metric1",
			MType: types.Gauge,
			Value: &metricValue,
		},
		{
			ID:    "metric2",
			MType: types.Counter,
			Value: &metricValue,
		},
	}

	t.Run("success returns HTML with metrics", func(t *testing.T) {
		mockSvc.EXPECT().
			List(gomock.Any()).
			Return(metrics, nil)

		handler := NewMetricListHTMLHandler(errHandlerFunc, mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "text/html; charset=utf-8", res.Header.Get("Content-Type"))

		body := rec.Body.String()
		assert.Contains(t, body, "<h1>Metrics</h1>")
		// Проверяем динамически отформатированные значения
		assert.Contains(t, body, "metric1: "+types.GetMetricStringValue(&metrics[0]))
		assert.Contains(t, body, "metric2: "+types.GetMetricStringValue(&metrics[1]))
	})

	t.Run("service error returns handled error", func(t *testing.T) {
		mockSvc.EXPECT().
			List(gomock.Any()).
			Return(nil, errors.New("some error"))

		handler := NewMetricListHTMLHandler(errHandlerFunc, mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		body := rec.Body.String()
		assert.True(t, strings.Contains(body, "some error"))
	})
}
