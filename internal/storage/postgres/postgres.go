package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Sweethear11/msg-processing-service/internal/config"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(log *slog.Logger, cfg *config.Config) (*Storage, error) {
	op := "storage.postgres.New"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	log.Debug("connection string", "connection string", connectionString)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateMessage(ctx context.Context, msg string) (int, error) {
	const op = "storage.postgres.createMessage"

	var id int
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO messages (message) VALUES ($1) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, msg).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetStats(ctx context.Context) (map[string]int, error) {
	const op = "storage.postgres.getStats"

	stats := make(map[string]int)

	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM messages WHERE processed = TRUE").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stats["processed"] = count

	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM messages WHERE processed = FALSE").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stats["unprocessed"] = count

	return stats, nil
}

func (s *Storage) MarkMessageAsProcessed(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE messages SET processed = TRUE WHERE id = $1", id)
	return err
}
