package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

func TestMetricUpdatePathHandler_TableDriven(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	tests := []struct {
		name            string
		args            args
		valFunc         func(string, string, string) error
		errValHandler   func(error) *types.APIError
		mockSvcBehavior func(m *MockMetricPathUpdater, metrics []types.Metrics)
		wantStatusCode  int
		wantBodyContain string
	}{
		{
			name: "OK",
			args: args{"gauge", "heap", "123.4"},
			valFunc: func(mt, mn, mv string) error {
				return nil
			},
			errValHandler: func(err error) *types.APIError {
				return nil
			},
			mockSvcBehavior: func(m *MockMetricPathUpdater, metrics []types.Metrics) {
				m.EXPECT().Update(gomock.Any(), gomock.Len(1)).Return(nil)
			},
			wantStatusCode:  http.StatusOK,
			wantBodyContain: "",
		},
		{
			name: "Validation error",
			args: args{"gauge", "heap", "not-a-float"},
			valFunc: func(mt, mn, mv string) error {
				return errors.New("bad metric")
			},
			errValHandler: func(err error) *types.APIError {
				return &types.APIError{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				}
			},
			mockSvcBehavior: func(m *MockMetricPathUpdater, metrics []types.Metrics) {
				// Service update should NOT be called
			},
			wantStatusCode:  http.StatusBadRequest,
			wantBodyContain: "bad metric",
		},
		{
			name: "Internal service error",
			args: args{"counter", "hits", "42"},
			valFunc: func(mt, mn, mv string) error {
				return nil
			},
			errValHandler: func(err error) *types.APIError {
				return nil
			},
			mockSvcBehavior: func(m *MockMetricPathUpdater, metrics []types.Metrics) {
				m.EXPECT().Update(gomock.Any(), gomock.Len(1)).Return(errors.New("DB failure"))
			},
			wantStatusCode:  http.StatusInternalServerError,
			wantBodyContain: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := NewMockMetricPathUpdater(ctrl)
			if tt.mockSvcBehavior != nil {
				tt.mockSvcBehavior(mockSvc, nil)
			}

			url := "/update/" + tt.args.metricType + "/" + tt.args.metricName + "/" + tt.args.metricValue
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(nil))
			req = addChiParams(req, "type", tt.args.metricType, "name", tt.args.metricName, "value", tt.args.metricValue)
			rec := httptest.NewRecorder()

			handler := NewMetricUpdatePathHandler(tt.valFunc, tt.errValHandler, mockSvc)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatusCode, rec.Code)
			if tt.wantBodyContain != "" {
				assert.Contains(t, rec.Body.String(), tt.wantBodyContain)
			}
		})
	}
}

// Helper to add URL params to chi context
func addChiParams(req *http.Request, params ...string) *http.Request {
	rctx := chi.NewRouteContext()
	for i := 0; i < len(params); i += 2 {
		rctx.URLParams.Add(params[i], params[i+1])
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}
