package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"wallet_controller/cmd/app"
)

func main() {
	slog.Info("Starting main")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error)

	go func() {
		if err := app.StartApplication(ctx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("Stopping main")

	case err := <-errChan:
		if err != nil {
			slog.Error("Error during application start:", err.Error(), nil)
		} else {
			slog.Info("Application stopped")
		}
	}

	slog.Info("Shutdown completed")
}
