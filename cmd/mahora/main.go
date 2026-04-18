package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancle := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGABRT,
	)
	defer cancle()

	slog.Info("mahora Starting")
	<-ctx.Done()
	slog.Info("shtting down")
}
