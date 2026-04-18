package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pavelc4/mahora/internal/auth"
	"github.com/pavelc4/mahora/internal/bot"
	"github.com/pavelc4/mahora/internal/config"
	"github.com/pavelc4/mahora/internal/db"
	"github.com/pavelc4/mahora/internal/github"
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
	log.Info("mahora starting", "env", cfg.Env)

	database, err := db.New(ctx, cfg.DBPath)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}
	defer database.Close()

	poller := auth.NewPoller(cfg.WorkerURL, cfg.WorkerSecret)

	b, err := bot.New(cfg, database.Queries, poller)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	githubPoller := github.NewPoller(database.Queries, b.Tele(), 5*time.Minute)

	go b.Start()
	go githubPoller.Run(ctx)

	log.Info("mahora started")
	<-ctx.Done()

	log.Info("shutting down")
	b.Stop()
	return nil
}
