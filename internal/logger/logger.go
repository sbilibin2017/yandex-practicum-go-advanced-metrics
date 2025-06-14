package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger = zap.NewNop().Sugar()

// Initialize configures the global logger with the specified log level.
//
// The level parameter should be a string recognized by zap.ParseAtomicLevel
// (e.g., "debug", "info", "warn", "error").
//
// Returns an error if the level string is invalid.
//
// On success, replaces the package-level Log variable with a SugaredLogger
// configured to the given log level.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	baseLogger, _ := cfg.Build()
	Log = baseLogger.Sugar()
	return nil
}
