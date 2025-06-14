package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags_Defaults(t *testing.T) {
	os.Clearenv()
	os.Args = []string{"cmd"}

	cfg, err := parseFlags()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:8080", cfg.ServerAddress)
	assert.Equal(t, "/update", cfg.ServerEndpoint)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 2, cfg.PollInterval)
	assert.Equal(t, 10, cfg.ReportInterval)
	assert.Equal(t, 5, cfg.NumWorkers)
}

func TestParseFlags_EnvOverrides(t *testing.T) {
	t.Setenv("SERVER_ADDRESS", "env:1234")
	t.Setenv("SERVER_ENDPOINT", "/env-update")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("POLL_INTERVAL", "99")
	t.Setenv("REPORT_INTERVAL", "100")
	t.Setenv("NUM_WORKERS", "7")

	os.Args = []string{"cmd", "-a=flag:5678", "-e=/flag-update", "-l=warn", "-p=1", "-r=2", "-w=3"}

	cfg, err := parseFlags()
	assert.NoError(t, err)

	assert.Equal(t, "env:1234", cfg.ServerAddress)
	assert.Equal(t, "/env-update", cfg.ServerEndpoint)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 99, cfg.PollInterval)
	assert.Equal(t, 100, cfg.ReportInterval)
	assert.Equal(t, 7, cfg.NumWorkers)
}

func TestParseFlags_FlagsFallback(t *testing.T) {
	os.Clearenv()
	os.Args = []string{
		"cmd",
		"-a=flaghost:9999",
		"-e=/metrics",
		"-l=trace",
		"-p=11",
		"-r=12",
		"-w=13",
	}

	cfg, err := parseFlags()
	assert.NoError(t, err)

	assert.Equal(t, "flaghost:9999", cfg.ServerAddress)
	assert.Equal(t, "/metrics", cfg.ServerEndpoint)
	assert.Equal(t, "trace", cfg.LogLevel)
	assert.Equal(t, 11, cfg.PollInterval)
	assert.Equal(t, 12, cfg.ReportInterval)
	assert.Equal(t, 13, cfg.NumWorkers)
}

func TestParseFlags_InvalidEnvFallbacksToFlag(t *testing.T) {
	t.Setenv("POLL_INTERVAL", "badint")
	t.Setenv("REPORT_INTERVAL", "badint")
	t.Setenv("NUM_WORKERS", "badint")

	os.Args = []string{
		"cmd",
		"-p=21",
		"-r=22",
		"-w=23",
	}

	cfg, err := parseFlags()
	assert.NoError(t, err)

	assert.Equal(t, 21, cfg.PollInterval)
	assert.Equal(t, 22, cfg.ReportInterval)
	assert.Equal(t, 23, cfg.NumWorkers)
}
