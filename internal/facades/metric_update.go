package facades

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/types"
)

type MetricUpdateFacade struct {
	client     *resty.Client
	serverAddr string
	endpoint   string
}

func NewMetricUpdateFacade(client *resty.Client, serverAddr string, endpoint string) *MetricUpdateFacade {
	if !strings.HasPrefix(serverAddr, "http://") && !strings.HasPrefix(serverAddr, "https://") {
		serverAddr = "http://" + serverAddr
	}
	return &MetricUpdateFacade{
		client:     client,
		serverAddr: serverAddr,
		endpoint:   endpoint,
	}
}

func (f *MetricUpdateFacade) Update(ctx context.Context, metrics types.Metrics) error {
	// Формируем URL без лишних слешей
	url := fmt.Sprintf("%s/%s", f.serverAddr, f.endpoint)

	// Отладочный лог: сериализация метрики в JSON
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshal metrics error: %w", err)
	}
	fmt.Printf("Sending POST to %s with body: %s\n", url, string(jsonData))

	resp, err := f.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
		Post(url)

	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
