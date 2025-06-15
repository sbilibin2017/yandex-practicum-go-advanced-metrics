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
	var addrFlag string
	fs.StringVar(&addrFlag, "a", "localhost:8080", "HTTP server endpoint address")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("ADDRESS"); env != "" {
			cfg.ServerAddress = env
			return
		}
		cfg.ServerAddress = addrFlag
	}
}

func withServerEndpoint(fs *flag.FlagSet) configs.AgentOption {
	var endpointFlag string
	fs.StringVar(&endpointFlag, "e", "/update", "API endpoint for updating metrics")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("SERVER_ENDPOINT"); env != "" {
			cfg.ServerEndpoint = env
			return
		}
		cfg.ServerEndpoint = endpointFlag
	}
}

func withLogLevel(fs *flag.FlagSet) configs.AgentOption {
	var levelFlag string
	fs.StringVar(&levelFlag, "l", "info", "logging level")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("LOG_LEVEL"); env != "" {
			cfg.LogLevel = env
			return
		}
		cfg.LogLevel = levelFlag
	}
}

func withPollInterval(fs *flag.FlagSet) configs.AgentOption {
	var pollFlag int
	fs.IntVar(&pollFlag, "p", 2, "metric polling frequency (pollInterval) in seconds")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("POLL_INTERVAL"); env != "" {
			if v, err := strconv.Atoi(env); err == nil && v > 0 {
				cfg.PollInterval = v
				return
			}
		}
		cfg.PollInterval = pollFlag
	}
}

func withReportInterval(fs *flag.FlagSet) configs.AgentOption {
	var reportFlag int
	fs.IntVar(&reportFlag, "r", 10, "metric reporting frequency (reportInterval) in seconds")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("REPORT_INTERVAL"); env != "" {
			if v, err := strconv.Atoi(env); err == nil && v > 0 {
				cfg.ReportInterval = v
				return
			}
		}
		cfg.ReportInterval = reportFlag
	}
}

func withNumWorkers(fs *flag.FlagSet) configs.AgentOption {
	var workersFlag int
	fs.IntVar(&workersFlag, "w", 5, "number of workers")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("NUM_WORKERS"); env != "" {
			if v, err := strconv.Atoi(env); err == nil && v > 0 {
				cfg.NumWorkers = v
				return
			}
		}
		cfg.NumWorkers = workersFlag
	}
}
