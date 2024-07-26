package kafka

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/Sweethear11/msg-processing-service/internal/storage/postgres"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	store      *postgres.Storage
	log        *slog.Logger
	broker     string
	topic      string
	maxWorkers int
}

func NewConsumer(store *postgres.Storage, log *slog.Logger, broker, topic string, maxWorkers int) *Consumer {
	return &Consumer{
		store:      store,
		log:        log,
		broker:     broker,
		topic:      topic,
		maxWorkers: maxWorkers,
	}
}

func (c *Consumer) ConsumeMessages(ctx context.Context) {
	op := "kafka.ConsumeMessages"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.broker},
		Topic:   c.topic,
	})
	defer reader.Close()

	sem := make(chan struct{}, c.maxWorkers)
	var wg sync.WaitGroup

	for {
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				c.log.Info("could not read message", "op", op, "err", err)
				return
			}

			c.log.Info("consumer received message", "message", string(msg.Value))

			id, err := strconv.Atoi(string(msg.Key))
			if err != nil {
				c.log.Info("could not parse message key", "op", op, "err", err)
				return
			}

			// Imitate processing time
			time.Sleep(5 * time.Second)

			c.store.MarkMessageAsProcessed(ctx, id)

			c.log.Info("new message marked as processed", slog.String("message", string(msg.Value)))
		}()

	}
}
