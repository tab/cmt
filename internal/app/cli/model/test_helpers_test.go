package model

import (
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func newTestLogger() logger.Logger {
	cfg := config.DefaultConfig()
	cfg.Logging.Format = logger.JSONFormat
	cfg.Logging.Level = logger.InfoLevel
	return logger.NewLogger(cfg)
}
