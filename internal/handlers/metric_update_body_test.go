package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricUpdateBodyHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricBodyUpdater(ctrl)

	validate := func(m types.Metrics) error {
		if m.ID == "invalid" {
			return errors.New("validation failed")
		}
		return nil
	}

	errorHandler := func(err error) *types.APIError {
		if err == nil {
			return nil
		}
		if err.Error() == "validation failed" {
			return &types.APIError{
				Code:    http.StatusBadRequest,
				Message: "validation error",
			}
		}
		return &types.APIError{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
		}
	}

	handler := NewMetricUpdateBodyHandler(validate, errorHandler, mockSvc)

	tests := []struct {
		name               string
		contentType        string
		body               string
		setupMock          func()
		wantStatus         int
		wantResponseSubstr string
	}{
		{
			name:        "invalid content-type",
			contentType: "text/plain",
			body:        `{"id":"metric1","mtype":"gauge","value":10}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid json body",
			contentType: "application/json",
			body:        `{invalid json}`,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:               "validation error",
			contentType:        "application/json",
			body:               `{"id":"invalid","mtype":"gauge","value":10}`,
			wantStatus:         http.StatusBadRequest,
			wantResponseSubstr: "validation error",
		},
		{
			name:        "service update error",
			contentType: "application/json",
			body:        `{"id":"metric1","mtype":"gauge","value":10}`,
			setupMock: func() {
				mockSvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("service failure"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:        "successful update",
			contentType: "application/json",
			body:        `{"id":"metric1","mtype":"gauge","value":10}`,
			setupMock: func() {
				mockSvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantStatus:         http.StatusOK,
			wantResponseSubstr: `"id":"metric1"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()

			handler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantResponseSubstr != "" {
				assert.Contains(t, rec.Body.String(), tt.wantResponseSubstr)
			}
		})
	}
}

// errorWriter implements http.ResponseWriter but returns error on Write
type errorWriter struct {
	header http.Header
}

func (e *errorWriter) Header() http.Header {
	return e.header
}

func (e *errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func (e *errorWriter) WriteHeader(statusCode int) {}

func TestNewMetricUpdateBodyHandler_EncodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricBodyUpdater(ctrl)

	validate := func(m types.Metrics) error {
		return nil // valid metric
	}

	errorHandler := func(err error) *types.APIError {
		return nil // no validation error
	}

	handler := NewMetricUpdateBodyHandler(validate, errorHandler, mockSvc)

	mockSvc.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(`{"id":"metric1","mtype":"gauge","value":10}`))
	req.Header.Set("Content-Type", "application/json")

	w := &errorWriter{header: make(http.Header)}

	handler(w, req)

}
