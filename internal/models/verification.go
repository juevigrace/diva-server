package models

import (
	"time"
)

type UserVerification struct {
	UserAction *UserAction
	Token      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}
