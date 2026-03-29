package models

import (
	"time"

	"github.com/google/uuid"
)

type UserVerification struct {
	UserID     uuid.UUID
	UserAction *UserAction
	Token      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}
