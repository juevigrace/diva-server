package models

import "github.com/google/uuid"

type UserPreferences struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Theme               Theme
	OnboardingCompleted bool
	Language            string
	LastSyncAt          int64
	CreatedAt           int64
	UpdatedAt           int64
}
