package model

import "time"

type Message struct {
	ID        int       `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	Processed bool      `json:"processed" db:"processed"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
