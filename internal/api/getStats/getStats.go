package getstats

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Total             int `json:"total"`
	ProcessedMessages int `json:"processedMessages"`
}
type MessageService interface {
	GetStats(ctx context.Context) (map[string]int, error)
}

func GetStats(log *slog.Logger, svc MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.handler.GetStats"

		w.Header().Set("Content-Type", "application/json")

		stats, err := svc.GetStats(r.Context())
		if err != nil {
			log.Error("failed to fetch statistics", "op", op, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("failed to fetch statistics")
			return
		}

		res := Response{
			Total:             stats["total"],
			ProcessedMessages: stats["processed"],
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error("failed to encode response", "op", op, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("failed to fetch statistics")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
