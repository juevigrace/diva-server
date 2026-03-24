package models

import "github.com/google/uuid"

type UserAction struct {
	Action Action
	UserID uuid.UUID
}
