package models

import (
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

type Session struct {
	ID           uuid.UUID
	User         User
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

func ToSessionResponse(s *Session) *responses.SessionResponse {
	return &responses.SessionResponse{
		SessionId:    s.ID.String(),
		UserId:       s.User.ID.String(),
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		Status:       s.Status.String(),
		Device:       s.Device,
		Ip:           s.IpAddress,
		Agent:        s.UserAgent,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
