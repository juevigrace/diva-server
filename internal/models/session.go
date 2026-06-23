package models

import (
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage"
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
	Type         SessionType
	ExpiresAt    int64
	CreatedAt    int64
	UpdatedAt    int64
}

func (s *Session) Response() *responses.SessionResponse {
	return &responses.SessionResponse{
		SessionId:    s.ID.String(),
		UserId:       s.User.ID.String(),
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		Status:       s.Status.String(),
		Type:         s.Type.String(),
		Device:       s.Device,
		Ip:           s.IpAddress,
		Agent:        s.UserAgent,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

func (s *Session) DBCreate() *storage.CreateSessionParams {
	return &storage.CreateSessionParams{
		ID:           s.ID,
		UserID:       s.User.ID,
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		Status:       s.Status.ToDB(),
		Type:         s.Type.ToDB(),
		Device:       s.Device,
		IpAddress:    s.IpAddress,
		UserAgent:    s.UserAgent,
		ExpiresAt:    s.ExpiresAt,
	}
}

func (s *Session) DBUpdate() *storage.UpdateSessionParams {
	return &storage.UpdateSessionParams{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		IpAddress:    s.IpAddress,
		ExpiresAt:    s.ExpiresAt,
		ID:           s.ID,
	}
}

func SessionFromDB(row *storage.DivaSession) *Session {
	return &Session{
		ID:           row.ID,
		User:         User{ID: row.UserID},
		AccessToken:  row.AccessToken,
		RefreshToken: row.RefreshToken,
		Device:       row.Device,
		IpAddress:    row.IpAddress,
		UserAgent:    row.UserAgent,
		Status:       SessionStatusFromDB(row.Status),
		Type:         SessionTypeFromDB(row.Type),
		ExpiresAt:    row.ExpiresAt,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
