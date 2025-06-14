package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/yandex-practicum-go-advanced-metrics/internal/configs"
)

// parseFlags parses command-line flags and environment variables to create a ServerConfig.
//
// It defines flags for server address (-a) and log level (-l), with default values.
// Environment variables ADDRESS and LOG_LEVEL override the flag values if set.
// Returns the constructed ServerConfig and any parsing error (currently always nil).
func parseFlags() (*configs.ServerConfig, error) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)

	options := []configs.ServerOption{
		withAddr(fs),
		withLogLevel(fs),
	}

	fs.Parse(os.Args[1:])

	return configs.NewServerConfig(options...), nil
}

// withAddr returns a ServerOption that sets the server address.
//
// It defines a flag "-a" for address with default ":8080".
// If the environment variable ADDRESS is set, it takes precedence over the flag.
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

// withLogLevel returns a ServerOption that sets the log level.
//
// It defines a flag "-l" for log level with default "info".
// If the environment variable LOG_LEVEL is set, it overrides the flag value.
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
