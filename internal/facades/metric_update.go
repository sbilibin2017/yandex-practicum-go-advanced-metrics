package facades

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

// MetricUpdateFacade provides a facade for sending HTTP requests
// to update metric values on a metrics server.
type MetricUpdateFacade struct {
	client     *resty.Client // HTTP client used for making requests
	serverAddr string        // Metrics server address
	endpoint   string        // Base endpoint path for updating metrics
}

// NewMetricUpdateFacade creates and returns a new instance of MetricUpdateFacade,
//
// Parameters:
//   - client: a configured instance of resty.Client
//   - serverAddr: address of the metrics server (e.g., "localhost:8080")
//   - endpoint: endpoint path for metric updates (e.g., "update")
//
// Returns:
//   - *MetricUpdateFacade: an initialized facade for sending metric updates.
func NewMetricUpdateFacade(client *resty.Client, serverAddr string, endpoint string) *MetricUpdateFacade {
	return &MetricUpdateFacade{
		client:     client,
		serverAddr: strings.TrimRight(serverAddr, "/"),
		endpoint:   strings.Trim(strings.TrimLeft(endpoint, "/"), "/"),
	}
}

// Update constructs a URL and sends a POST request to update a metric.
//
// The URL is constructed using the pattern: /{endpoint}/{type}/{name}/{value}.
// Example: /update/gauge/Alloc/123.45
//
// Parameters:
//   - ctx: context for request cancellation and timeout
//   - req: a MetricsUpdatePathRequest containing the metric name, type, and value
//
// Returns:
//   - error: if the request fails or the server responds with a bad status code.
func (f *MetricUpdateFacade) Update(ctx context.Context, req types.MetricsUpdatePathRequest) error {
	addr := f.serverAddr
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}

	url := fmt.Sprintf("%s/%s/%s/%s/%s", addr, f.endpoint, req.MType, req.Name, req.Value)

	resp, err := f.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
