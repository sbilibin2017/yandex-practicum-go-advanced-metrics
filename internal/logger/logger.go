package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger = zap.NewNop().Sugar()

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
