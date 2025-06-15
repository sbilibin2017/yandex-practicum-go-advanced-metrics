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
	"github.com/stretchr/testify/assert"
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

	// Start httptest server with the app's handler
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
		{"Valid gauge metric", "gauge", "temperature", "42.5", http.StatusOK},
		{"Gauge override previous value", "gauge", "temperature", "24.3", http.StatusOK},
		{"Valid counter metric", "counter", "requests", "100", http.StatusOK},
		{"Counter accumulate value", "counter", "requests", "50", http.StatusOK},
		{"Invalid metric type", "invalid", "some", "10", http.StatusBadRequest},
		{"Invalid metric value", "gauge", "pressure", "not-a-number", http.StatusBadRequest},
		{"Missing metric value", "gauge", "humidity", "", http.StatusBadRequest},
		// Missing metric name or type should return 404 as per your server routing
		{"Missing metric name", "gauge", "", "10.1", http.StatusNotFound},
		{"Missing metric type", "", "foo", "1.1", http.StatusNotFound},
		{"Empty path", "", "", "", http.StatusNotFound},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var url string
			if tt.metricType != "" && tt.metricName != "" && tt.metricValue != "" {
				url = fmt.Sprintf("/update/%s/%s/%s", tt.metricType, tt.metricName, tt.metricValue)
			} else if tt.metricType != "" && tt.metricName != "" {
				// This path does not exist for POST without value, so expect 404
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

			// Check request error before any assertions
			s.Require().NoError(err, "HTTP error in request for URL: %s", url)
			s.Equal(tt.expectedStatus, resp.StatusCode(), "Unexpected status for %s", url)
		})
	}
}

func (s *UpdateMetricSuite) TestMetricValuePathHandler() {
	// Pre-load some metrics to get their values later
	resp, err := s.client.R().
		SetContext(context.Background()).
		SetHeader("Content-Type", "text/plain").
		Post("/update/counter/requests/150")
	s.Require().NoError(err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode())

	resp, err = s.client.R().
		SetContext(context.Background()).
		SetHeader("Content-Type", "text/plain").
		Post("/update/gauge/temperature/42.5")
	s.Require().NoError(err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode())

	tests := []struct {
		name           string
		metricType     string
		metricName     string
		expectedStatus int
	}{
		{"Get gauge value", "gauge", "temperature", http.StatusOK},
		{"Get counter value", "counter", "requests", http.StatusOK},
		// {"Unknown metric type", "invalid", "foo", http.StatusBadRequest},
		// {"Missing metric name", "gauge", "", http.StatusNotFound},
		// {"Missing metric type", "", "foo", http.StatusBadRequest},
		// {"Empty path", "", "", http.StatusNotFound},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			url := fmt.Sprintf("/value/%s/%s", tt.metricType, tt.metricName)
			resp, err := s.client.R().
				SetContext(context.Background()).
				Get(url)

			s.Require().NoError(err, "HTTP error in request for URL: %s", url)
			s.Equal(tt.expectedStatus, resp.StatusCode(), "Unexpected status for %s", url)
		})
	}
}

func (s *UpdateMetricSuite) TestMetricsListHandler() {
	resp, err := s.client.R().
		SetContext(context.Background()).
		Get("/")

	s.Require().NoError(err, "HTTP error in request for /")
	s.Equal(http.StatusOK, resp.StatusCode(), "Expected 200 OK from metrics list handler")
	s.Contains(resp.Header().Get("Content-Type"), "text/html", "Expected Content-Type text/html for list")
	s.Contains(resp.String(), "<html>", "Expected HTML response body")
}

func TestUpdateMetricSuite(t *testing.T) {
	suite.Run(t, new(UpdateMetricSuite))
}
