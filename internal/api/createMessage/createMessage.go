package createmessage

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Sweethear11/msg-processing-service/internal/model"
)

type Request struct {
	Msg string `json:"message"`
}

type MessageService interface {
	CreateMessage(ctx context.Context, msg string) (model.Message, error)
}

func NewMessage(log *slog.Logger, svc MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.handler.NewMessage"

		w.Header().Set("Content-Type", "application/json")

		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", "op", op, "err", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("invalid request")
			return
		}

		if req.Msg == "" {
			log.Error("message is an empty string")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("empty message")
			return
		}

		message, err := svc.CreateMessage(r.Context(), req.Msg)
		if err != nil {
			log.Error("failed to create a message", "op", op, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("failed to create message")
			return
		}

		log.Info("new message", slog.String("message", req.Msg))
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(message)
	}
}
