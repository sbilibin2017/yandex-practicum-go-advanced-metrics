package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
)

func parseFlags() (*configs.ServerConfig, error) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)

	options := []configs.ServerOption{
		withAddr(fs),
		withLogLevel(fs),
	}

	fs.Parse(os.Args[1:])

	return configs.NewServerConfig(options...), nil
}

func withAddr(fs *flag.FlagSet) configs.ServerOption {
	var addr string
	fs.StringVar(&addr, "a", ":8080", "address and port to run server")

	return func(cfg *configs.ServerConfig) {
		if env := os.Getenv("ADDRESS"); env != "" {
			cfg.Address = env
		} else {
			cfg.Address = addr
		}
	}
}

func withLogLevel(fs *flag.FlagSet) configs.ServerOption {
	var level string
	fs.StringVar(&level, "l", "info", "log level")

	return func(cfg *configs.ServerConfig) {
		if env := os.Getenv("LOG_LEVEL"); env != "" {
			cfg.LogLevel = env
		} else {
			cfg.LogLevel = level
		}
	}
}
