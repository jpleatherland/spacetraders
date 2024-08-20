// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

type Session struct {
	ID        string        `json:"id"`
	ExpiresAt time.Time     `json:"expires_at"`
	UserID    uuid.UUID     `json:"user_id"`
	AgentID   uuid.NullUUID `json:"agent_id"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
}
