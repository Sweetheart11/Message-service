package model

import "time"

type Message struct {
	ID        int       `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	processed bool      `json:"processed" db:"processed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
