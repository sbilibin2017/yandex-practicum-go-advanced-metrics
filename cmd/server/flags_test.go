package main

import (
	"flag"
	"os"
	"testing"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	origArgs := os.Args
	origAddr := os.Getenv("ADDRESS")
	origLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		os.Args = origArgs
		os.Setenv("ADDRESS", origAddr)
		os.Setenv("LOG_LEVEL", origLogLevel)
	}()

	tests := []struct {
		name        string
		args        []string
		envAddress  string
		envLogLevel string
		wantAddr    string
		wantLevel   string
	}{
		{
			name:      "default flags no env",
			args:      []string{"cmd"},
			wantAddr:  ":8080",
			wantLevel: "info",
		},
		{
			name:      "flags override defaults",
			args:      []string{"cmd", "-a", ":9090", "-l", "debug"},
			wantAddr:  ":9090",
			wantLevel: "debug",
		},
		{
			name:        "env overrides flags",
			args:        []string{"cmd", "-a", ":9090", "-l", "debug"},
			envAddress:  ":7070",
			envLogLevel: "warn",
			wantAddr:    ":7070",
			wantLevel:   "warn",
		},
		{
			name:        "env only no flags",
			args:        []string{"cmd"},
			envAddress:  ":6060",
			envLogLevel: "error",
			wantAddr:    ":6060",
			wantLevel:   "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			if tt.envAddress != "" {
				os.Setenv("ADDRESS", tt.envAddress)
			} else {
				os.Unsetenv("ADDRESS")
			}

			if tt.envLogLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envLogLevel)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			cfg, err := parseFlags()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantAddr, cfg.Address)
			assert.Equal(t, tt.wantLevel, cfg.LogLevel)
		})
	}
}

func TestWithAddr(t *testing.T) {
	tests := []struct {
		name     string
		flagArgs []string
		envAddr  string
		wantAddr string
	}{
		{
			name:     "flag only",
			flagArgs: []string{"-a", ":6060"},
			wantAddr: ":6060",
		},
		{
			name:     "env overrides flag",
			flagArgs: []string{"-a", ":6060"},
			envAddr:  ":5050",
			wantAddr: ":5050",
		},
		{
			name:     "env only no flag",
			flagArgs: []string{},
			envAddr:  ":4040",
			wantAddr: ":4040",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envAddr != "" {
				os.Setenv("ADDRESS", tt.envAddr)
			} else {
				os.Unsetenv("ADDRESS")
			}

			fs := flag.NewFlagSet("test", flag.ExitOnError)
			opt := withAddr(fs)
			fs.Parse(tt.flagArgs)

			cfg := &configs.ServerConfig{}
			opt(cfg)
			assert.Equal(t, tt.wantAddr, cfg.Address)
		})
	}
}

func TestWithLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		flagArgs  []string
		envLevel  string
		wantLevel string
	}{
		{
			name:      "flag only",
			flagArgs:  []string{"-l", "error"},
			wantLevel: "error",
		},
		{
			name:      "env overrides flag",
			flagArgs:  []string{"-l", "error"},
			envLevel:  "fatal",
			wantLevel: "fatal",
		},
		{
			name:      "env only no flag",
			flagArgs:  []string{},
			envLevel:  "warn",
			wantLevel: "warn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envLevel)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			fs := flag.NewFlagSet("test", flag.ExitOnError)
			opt := withLogLevel(fs)
			fs.Parse(tt.flagArgs)

			cfg := &configs.ServerConfig{}
			opt(cfg)
			assert.Equal(t, tt.wantLevel, cfg.LogLevel)
		})
	}
}
