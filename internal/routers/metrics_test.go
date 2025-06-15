package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsRouter_HandlersAndMiddleware(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		url                 string
		expectStatus        int
		expectMiddleware    bool
		expectUpdateHandler bool
		expectValueHandler  bool
		expectListHandler   bool
	}{
		{
			name:                "POST /update route",
			method:              "POST",
			url:                 "/update/counter/testmetric/123",
			expectStatus:        http.StatusOK,
			expectMiddleware:    true,
			expectUpdateHandler: true,
		},
		{
			name:               "GET /value route",
			method:             "GET",
			url:                "/value/gauge/testmetric",
			expectStatus:       http.StatusOK,
			expectMiddleware:   true,
			expectValueHandler: true,
		},
		{
			name:              "GET / route",
			method:            "GET",
			url:               "/",
			expectStatus:      http.StatusOK,
			expectMiddleware:  true,
			expectListHandler: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var middlewareCalled, updateHandlerCalled, valueHandlerCalled, listHandlerCalled bool

			middleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					middlewareCalled = true
					next.ServeHTTP(w, r)
				})
			}

			updateHandler := func(w http.ResponseWriter, r *http.Request) {
				updateHandlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("update-ok"))
			}

			valueHandler := func(w http.ResponseWriter, r *http.Request) {
				valueHandlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("value-ok"))
			}

			listHandler := func(w http.ResponseWriter, r *http.Request) {
				listHandlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("list-ok"))
			}

			router := NewMetricsRouter(updateHandler, valueHandler, listHandler, middleware)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.expectStatus, resp.StatusCode)
			assert.Equal(t, tt.expectMiddleware, middlewareCalled, "middleware called")
			assert.Equal(t, tt.expectUpdateHandler, updateHandlerCalled, "updateHandler called")
			assert.Equal(t, tt.expectValueHandler, valueHandlerCalled, "valueHandler called")
			assert.Equal(t, tt.expectListHandler, listHandlerCalled, "listHandler called")
		})
	}
}
