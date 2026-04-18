package config

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(cfg *Config) *slog.Logger {
	return slog.New(newHandler(cfg))
}

func newHandler(cfg *Config) slog.Handler {
	opts := &slog.HandlerOptions{Level: parseLevel(cfg.LogLevel)}

	if cfg.Env == "production" {
		return slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.NewTextHandler(os.Stdout, opts)
}

var levelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func parseLevel(s string) slog.Level {
	if lvl, ok := levelMap[strings.ToLower(s)]; ok {
		return lvl
	}
	return slog.LevelInfo
}
