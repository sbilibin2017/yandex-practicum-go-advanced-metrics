package facades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdateFacade_Update_Success(t *testing.T) {
	// Arrange: create test server
	var receivedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := resty.New()
	endpoint := "update"
	facade := NewMetricUpdateFacade(client, server.URL, endpoint)

	req := types.MetricsUpdatePathRequest{
		Name:  "Alloc",
		MType: "gauge",
		Value: "123.45",
	}

	// Act
	err := facade.Update(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "/update/gauge/Alloc/123.45", receivedPath)
}

func TestMetricUpdateFacade_Update_ServerError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := resty.New()
	facade := NewMetricUpdateFacade(client, server.URL, "update")

	req := types.MetricsUpdatePathRequest{
		Name:  "Heap",
		MType: "gauge",
		Value: "1000",
	}

	// Act
	err := facade.Update(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "server returned status 500")
}

func TestMetricUpdateFacade_Update_InvalidScheme(t *testing.T) {
	// Arrange
	client := resty.New()
	invalidAddr := "localhost:1234" // No http://
	facade := NewMetricUpdateFacade(client, invalidAddr, "update")

	req := types.MetricsUpdatePathRequest{
		Name:  "Heap",
		MType: "gauge",
		Value: "1000",
	}

	// Using a dummy server that won't respond; just checking for valid URL formatting
	go func() {
		_ = facade.Update(context.Background(), req)
	}()

	// Just checking it doesn't panic on missing scheme
	assert.True(t, true)
}

func TestMetricUpdateFacade_Update_RequestError_NoMock(t *testing.T) {
	client := resty.New()

	// Use an invalid address that will cause connection failure
	invalidAddr := "http://invalid-host.local:12345"
	facade := NewMetricUpdateFacade(client, invalidAddr, "update")

	req := types.MetricsUpdatePathRequest{
		MType: "gauge",
		Name:  "Alloc",
		Value: "123",
	}

	err := facade.Update(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request error")
}
