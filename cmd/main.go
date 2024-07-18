package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Sweethear11/msg-processing-service/internal/config"
	"github.com/Sweethear11/msg-processing-service/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting time tracker service",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	store, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("cannot create storage", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", store)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
