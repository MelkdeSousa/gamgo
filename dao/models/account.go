package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"passwordHash,omitempty"`
	Email        string    `json:"email"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	DeletedAt    time.Time `json:"deletedAt"`
}
