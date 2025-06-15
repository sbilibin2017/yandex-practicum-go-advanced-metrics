package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
)

func parseFlags() (*configs.AgentConfig, error) {
	fs := flag.NewFlagSet("agent", flag.ExitOnError)

	options := []configs.AgentOption{
		withServerAddress(fs),
		withServerEndpoint(fs),
		withLogLevel(fs),
		withPollInterval(fs),
		withReportInterval(fs),
		withNumWorkers(fs),
	}

	fs.Parse(os.Args[1:])

	return configs.NewAgentConfig(options...), nil
}

func withServerAddress(fs *flag.FlagSet) configs.AgentOption {
	var addr string
	fs.StringVar(&addr, "a", "localhost:8080", "HTTP server endpoint address")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("SERVER_ADDRESS"); env != "" {
			cfg.ServerAddress = env
		} else {
			cfg.ServerAddress = addr
		}
	}
}

func withServerEndpoint(fs *flag.FlagSet) configs.AgentOption {
	var endpoint string
	fs.StringVar(&endpoint, "e", "/update", "API endpoint for updating metrics")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("SERVER_ENDPOINT"); env != "" {
			cfg.ServerEndpoint = env
		} else {
			cfg.ServerEndpoint = endpoint
		}
	}
}

func withLogLevel(fs *flag.FlagSet) configs.AgentOption {
	var level string
	fs.StringVar(&level, "l", "info", "logging level")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("LOG_LEVEL"); env != "" {
			cfg.LogLevel = env
		} else {
			cfg.LogLevel = level
		}
	}
}

func withPollInterval(fs *flag.FlagSet) configs.AgentOption {
	var poll int
	fs.IntVar(&poll, "p", 2, "metric polling frequency (pollInterval) in seconds")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("POLL_INTERVAL"); env != "" {
			if v, err := strconv.Atoi(env); err == nil {
				cfg.PollInterval = v
				return
			}
		}
		cfg.PollInterval = poll
	}
}

func withReportInterval(fs *flag.FlagSet) configs.AgentOption {
	var report int
	fs.IntVar(&report, "r", 10, "metric reporting frequency (reportInterval) in seconds")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("REPORT_INTERVAL"); env != "" {
			if v, err := strconv.Atoi(env); err == nil {
				cfg.ReportInterval = v
				return
			}
		}
		cfg.ReportInterval = report
	}
}

func withNumWorkers(fs *flag.FlagSet) configs.AgentOption {
	var workers int
	fs.IntVar(&workers, "w", 5, "number of workers")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("NUM_WORKERS"); env != "" {
			if v, err := strconv.Atoi(env); err == nil {
				cfg.NumWorkers = v
				return
			}
		}
		cfg.NumWorkers = workers
	}
}
