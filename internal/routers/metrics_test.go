package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// dummy handler that writes status 200 OK
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestNewMetricsRouter(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
	}{
		{
			name:         "POST update with type, name, value",
			method:       http.MethodPost,
			url:          "/update/gauge/metric1/100",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST update with type and name",
			method:       http.MethodPost,
			url:          "/update/counter/metric2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST update with body",
			method:       http.MethodPost,
			url:          "/update/",
			expectedCode: http.StatusOK,
		},
		{
			name:         "GET value with type and name",
			method:       http.MethodGet,
			url:          "/value/gauge/metric1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "GET value with type only",
			method:       http.MethodGet,
			url:          "/value/counter",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST value with body",
			method:       http.MethodPost,
			url:          "/value/",
			expectedCode: http.StatusOK,
		},
		{
			name:         "GET metrics list root",
			method:       http.MethodGet,
			url:          "/",
			expectedCode: http.StatusOK,
		},
		{
			name:         "404 not found",
			method:       http.MethodGet,
			url:          "/notfound",
			expectedCode: http.StatusNotFound,
		},
	}

	// Create router with dummy handlers and no middlewares
	router := NewMetricsRouter(
		dummyHandler,
		dummyHandler,
		dummyHandler,
		dummyHandler,
		dummyHandler,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
