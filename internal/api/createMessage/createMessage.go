package createmessage

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Request struct {
	Msg string `json:"message"`
}

type MessageService interface {
	CreateMessage(ctx context.Context, msg string) error
}

func NewMessage(log *slog.Logger, svc MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.handler.NewMessage"

		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", "op", op, "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Msg == "" {
			log.Error("message is an empty string")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.CreateMessage(r.Context(), req.Msg); err != nil {
			log.Error("failed to create a message", "op", op, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("new message", slog.String("message", req.Msg))
		w.WriteHeader(http.StatusAccepted)
	}
}
