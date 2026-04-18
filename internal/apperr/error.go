package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pavelc4/mahora/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	log := config.NewLogger(cfg)
	log.Info("mahora starting", "env", cfg.Env, "log_level", cfg.LogLevel)

	// TODO: init db
	// TODO: start bot

	<-ctx.Done()
	log.Info("shutting down")
	return nil
}
