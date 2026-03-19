package models

import (
	"time"

	"github.com/google/uuid"
)

type UserVerification struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
