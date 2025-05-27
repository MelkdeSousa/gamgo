package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	IsActive     bool      `json:"-"`
	CreatedAt    time.Time `json:"-"`
	DeletedAt    time.Time `json:"-"`
}
