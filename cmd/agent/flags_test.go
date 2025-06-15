package main

import (
	"os"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/stretchr/testify/assert"
)

type parseFlagsTestCase struct {
	name            string
	env             map[string]string
	args            []string
	expected        configs.AgentConfig
	expectParseFail bool
}

func TestParseFlags_TableDriven(t *testing.T) {
	testCases := []parseFlagsTestCase{
		{
			name: "Defaults",
			env:  map[string]string{},
			args: []string{"cmd"},
			expected: configs.AgentConfig{
				ServerAddress:  ":8080",
				ServerEndpoint: "update/",
				LogLevel:       "info",
				PollInterval:   2,
				ReportInterval: 10,
				NumWorkers:     5,
			},
		},
		{
			name: "Env overrides flags",
			env: map[string]string{
				"ADDRESS":         "env:1234",
				"SERVER_ENDPOINT": "/env-update/",
				"LOG_LEVEL":       "debug",
				"POLL_INTERVAL":   "99",
				"REPORT_INTERVAL": "100",
				"NUM_WORKERS":     "7",
			},
			args: []string{"cmd", "-a=flag:5678", "-e=/flag-update/", "-l=warn", "-p=1", "-r=2", "-w=3"},
			expected: configs.AgentConfig{
				ServerAddress:  "env:1234",
				ServerEndpoint: "/env-update/",
				LogLevel:       "debug",
				PollInterval:   99,
				ReportInterval: 100,
				NumWorkers:     7,
			},
		},
		{
			name: "Flags used if no env",
			env:  map[string]string{},
			args: []string{"cmd", "-a=flaghost:9999", "-e=/metrics/", "-l=trace", "-p=11", "-r=12", "-w=13"},
			expected: configs.AgentConfig{
				ServerAddress:  "flaghost:9999",
				ServerEndpoint: "/metrics/",
				LogLevel:       "trace",
				PollInterval:   11,
				ReportInterval: 12,
				NumWorkers:     13,
			},
		},
		{
			name: "Invalid env fallback to flags",
			env: map[string]string{
				"POLL_INTERVAL":   "badint",
				"REPORT_INTERVAL": "badint",
				"NUM_WORKERS":     "badint",
			},
			args: []string{"cmd", "-a=localhost:8080", "-e=update/", "-l=info", "-p=21", "-r=22", "-w=23"},
			expected: configs.AgentConfig{
				ServerAddress:  "localhost:8080",
				ServerEndpoint: "update/",
				LogLevel:       "info",
				PollInterval:   21,
				ReportInterval: 22,
				NumWorkers:     23,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear env first
			os.Clearenv()
			// Set env vars from test case
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			// Set args
			os.Args = tc.args

			cfg, err := parseFlags()
			if tc.expectParseFail {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected.ServerAddress, cfg.ServerAddress)
			assert.Equal(t, tc.expected.ServerEndpoint, cfg.ServerEndpoint)
			assert.Equal(t, tc.expected.LogLevel, cfg.LogLevel)
			assert.Equal(t, tc.expected.PollInterval, cfg.PollInterval)
			assert.Equal(t, tc.expected.ReportInterval, cfg.ReportInterval)
			assert.Equal(t, tc.expected.NumWorkers, cfg.NumWorkers)
		})
	}
}
