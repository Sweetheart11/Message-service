package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	createmessage "github.com/Sweethear11/msg-processing-service/internal/api/createMessage"
	getstats "github.com/Sweethear11/msg-processing-service/internal/api/getStats"
	"github.com/Sweethear11/msg-processing-service/internal/config"
	kafkaService "github.com/Sweethear11/msg-processing-service/internal/kafka"
	"github.com/Sweethear11/msg-processing-service/internal/service"
	"github.com/Sweethear11/msg-processing-service/internal/storage/postgres"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
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
		"starting message procesing service",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	store, err := postgres.New(log, cfg)
	if err != nil {
		log.Error("cannot create storage", err)
		os.Exit(1)
	}

	go kafkaService.ConsumeMessages(context.Background(), store, log, "localhost:9092", "messages")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Broker),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	router.Post("/message", createmessage.NewMessage(log, service.NewMessageService(store, writer, log)))
	router.Get("/message", getstats.GetStats(log, service.NewMessageService(store, writer, log)))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.StringValue(err.Error()))

		return
	}

	log.Info("server stopped")

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
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
