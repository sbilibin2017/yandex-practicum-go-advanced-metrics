package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/apps"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/logger"
	"github.com/stretchr/testify/suite"
)

type UpdateMetricSuite struct {
	suite.Suite
	serverURL string
	client    *resty.Client
}

func (s *UpdateMetricSuite) SetupSuite() {
	config := &configs.ServerConfig{
		Address:  ":0",
		LogLevel: "debug",
	}

	err := logger.Initialize(config.LogLevel)
	s.Require().NoError(err)

	srv, err := apps.NewServerApp(config)
	s.Require().NoError(err)

	ts := httptest.NewServer(srv.Handler)
	s.T().Cleanup(ts.Close)

	s.serverURL = ts.URL
	s.client = resty.New().SetBaseURL(s.serverURL)
}

func (s *UpdateMetricSuite) TestUpdateMetricPathHandler() {
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		metricValue    string
		expectedStatus int
	}{
		{
			name:           "Valid gauge metric",
			metricType:     "gauge",
			metricName:     "temperature",
			metricValue:    "42.5",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Gauge override previous value",
			metricType:     "gauge",
			metricName:     "temperature",
			metricValue:    "24.3",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid counter metric",
			metricType:     "counter",
			metricName:     "requests",
			metricValue:    "100",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Counter accumulate value",
			metricType:     "counter",
			metricName:     "requests",
			metricValue:    "50",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid metric type",
			metricType:     "invalid",
			metricName:     "some",
			metricValue:    "10",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid metric value",
			metricType:     "gauge",
			metricName:     "pressure",
			metricValue:    "not-a-number",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing metric value",
			metricType:     "gauge",
			metricName:     "humidity",
			metricValue:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing metric name",
			metricType:     "gauge",
			metricName:     "",
			metricValue:    "10.1",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Missing metric type",
			metricType:     "",
			metricName:     "foo",
			metricValue:    "1.1",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty path",
			metricType:     "",
			metricName:     "",
			metricValue:    "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Construct path manually
			var url string
			if tt.metricType != "" && tt.metricName != "" && tt.metricValue != "" {
				url = fmt.Sprintf("/update/%s/%s/%s", tt.metricType, tt.metricName, tt.metricValue)
			} else if tt.metricType != "" && tt.metricName != "" {
				url = fmt.Sprintf("/update/%s/%s", tt.metricType, tt.metricName)
			} else if tt.metricType != "" {
				url = fmt.Sprintf("/update/%s", tt.metricType)
			} else {
				url = "/update"
			}

			resp, err := s.client.R().
				SetContext(context.Background()).
				SetHeader("Content-Type", "text/plain").
				Post(url)

			s.Require().NoError(err, "HTTP error in request")
			s.Equal(tt.expectedStatus, resp.StatusCode(), "Unexpected status for %s", url)
		})
	}
}

func TestUpdateMetricSuite(t *testing.T) {
	suite.Run(t, new(UpdateMetricSuite))
}
