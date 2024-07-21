package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage(ctx context.Context, log *slog.Logger, broker, topic, message string, id int) error {
	op := "kafka.ProduceMessage"
	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", id)),
		Value: []byte(message),
	})
	if err != nil {
		log.Info("could not write message", "op", op, "message", message)
		return err
	}
	return nil
}
