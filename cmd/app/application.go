package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"wallet_controller/config"
	"wallet_controller/internal/router"
	"wallet_controller/internal/storage"
)

func StartApplication(ctx context.Context) error {
	cfg := config.GetConfig()

	slog.Info("Starting application")

	cfg.Client = storage.NewConnection(ctx, cfg)

	r := router.SetupRouter(ctx, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Env.IpAddress, cfg.Env.ApiPort)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		slog.Info("Starting http server")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", err.Error())
			panic(err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down application")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", err.Error())
		panic(err)
	}
	return nil
}
