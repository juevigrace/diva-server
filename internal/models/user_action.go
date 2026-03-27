package models

import "github.com/google/uuid"

type UserAction struct {
	ID     uuid.UUID
	Action Action
	UserID uuid.UUID
}
