package models

import "github.com/google/uuid"

type Session struct {
	ID           uuid.UUID
	User         uuid.UUID
	AccessToken  string
	RefreshToken string
	Device       string
	IpAddress    string
	UserAgent    string
	Status       SessionStatus
	ExpiresAt    int64
	CreatedAt    int64
	UpdatedAt    int64
}
