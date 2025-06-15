package facades

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdateFacade_Update_Success(t *testing.T) {
	var receivedMetric types.Metrics
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		err := json.NewDecoder(r.Body).Decode(&receivedMetric)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := resty.New()
	facade := NewMetricUpdateFacade(client, server.URL, "update/")

	val := 123.45
	reqMetric := types.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &val,
	}

	err := facade.Update(context.Background(), reqMetric)
	require.NoError(t, err)
	assert.Equal(t, reqMetric, receivedMetric)
}

func TestMetricUpdateFacade_Update_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := resty.New()
	facade := NewMetricUpdateFacade(client, server.URL, "update/")

	val := 1000.0
	reqMetric := types.Metrics{
		ID:    "Heap",
		MType: "gauge",
		Value: &val,
	}

	err := facade.Update(context.Background(), reqMetric)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "server returned status 500")
}

func TestMetricUpdateFacade_Update_InvalidScheme(t *testing.T) {
	client := resty.New()
	invalidAddr := "localhost:1234" // no http:// prefix
	facade := NewMetricUpdateFacade(client, invalidAddr, "update")

	val := 1000.0
	reqMetric := types.Metrics{
		ID:    "Heap",
		MType: "gauge",
		Value: &val,
	}

	err := facade.Update(context.Background(), reqMetric)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request error")
}

func TestMetricUpdateFacade_Update_RequestError(t *testing.T) {
	client := resty.New()
	invalidAddr := "http://invalid-host.local:12345"
	facade := NewMetricUpdateFacade(client, invalidAddr, "update/")

	val := 123.0
	reqMetric := types.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &val,
	}

	err := facade.Update(context.Background(), reqMetric)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request error")
}
