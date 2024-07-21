package kafka

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/Sweethear11/msg-processing-service/internal/storage/postgres"
	"github.com/segmentio/kafka-go"
)

func ConsumeMessages(ctx context.Context, db *postgres.Storage, log *slog.Logger, broker, topic string) {
	op := "kafka.ConsumeMessages"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Info("could not read message: %v", "op", op, "err", err)
			continue
		}
		log.Info("received: %s", "message", string(msg.Value))

		id, err := strconv.Atoi(string(msg.Key))
		if err != nil {
			log.Info("could not parse message key: %v", "op", op, "err", err)
			continue
		}
		db.MarkMessageAsProcessed(ctx, id)

		log.Info("new message marked as processed", slog.String("message", string(msg.Value)))
	}
}
