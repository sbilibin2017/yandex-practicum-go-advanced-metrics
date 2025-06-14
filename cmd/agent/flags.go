package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
)

// parseFlags парсит флаги командной строки и переменные окружения,
// создавая AgentConfig с параметрами конфигурации агента.
//
// Используются флаги:
// -a=<ЗНАЧЕНИЕ> — адрес эндпоинта HTTP-сервера (по умолчанию "localhost:8080").
// -p=<ЗНАЧЕНИЕ> — pollInterval, частота опроса метрик (по умолчанию 2 секунды).
// -r=<ЗНАЧЕНИЕ> — reportInterval, частота отправки метрик на сервер (по умолчанию 10 секунд).
//
// Переменные окружения с приоритетом над флагами:
// SERVER_ADDRESS, POLL_INTERVAL, REPORT_INTERVAL.
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
	fs.StringVar(&addr, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера")

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
	fs.StringVar(&endpoint, "e", "/update", "API endpoint для обновления метрик")

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
	fs.StringVar(&level, "l", "info", "уровень логирования")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("LOG_LEVEL"); env != "" {
			cfg.LogLevel = env
		} else {
			cfg.LogLevel = level
		}
	}
}

// pollInterval по умолчанию 2 секунды
func withPollInterval(fs *flag.FlagSet) configs.AgentOption {
	var poll int
	fs.IntVar(&poll, "p", 2, "частота опроса метрик (pollInterval), сек")

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

// reportInterval по умолчанию 10 секунд
func withReportInterval(fs *flag.FlagSet) configs.AgentOption {
	var report int
	fs.IntVar(&report, "r", 10, "частота отправки метрик (reportInterval), сек")

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
	fs.IntVar(&workers, "w", 5, "количество воркеров")

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
