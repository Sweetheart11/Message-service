package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sweethear11/msg-processing-service/internal/model"
	storage "github.com/Sweethear11/msg-processing-service/internal/storage/postgres"
	"github.com/segmentio/kafka-go"
)

type MessageService struct {
	db       *storage.Storage
	producer *kafka.Writer
	log      *slog.Logger
}

func NewMessageService(db *storage.Storage, producer *kafka.Writer, log *slog.Logger) *MessageService {
	return &MessageService{
		db:       db,
		producer: producer,
		log:      log,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, message string) (model.Message, error) {
	op := "service.message.CreateMessage"
	id, err := s.db.CreateMessage(ctx, message)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.producer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", id)),
		Value: []byte(message),
	})
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.Message{
		ID:        id,
		Message:   message,
		CreatedAt: time.Now(),
	}, nil
}

func (s *MessageService) GetStats(ctx context.Context) (map[string]int, error) {
	op := "service.message.getStats"
	stats, err := s.db.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stats["total"] = stats["processed"] + stats["unprocessed"]
	s.log.Info("statistics fetched", slog.String("Total messages", fmt.Sprintf("%d", stats["processed"]+stats["unprocessed"])),
		slog.String("Processed messages", fmt.Sprintf("%d", stats["processed"])))

	return stats, nil
}
