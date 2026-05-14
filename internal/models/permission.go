package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	Action      string
	RoleLevel   Role
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   *int64
}

func (p *Permission) Response() *responses.PermissionResponse {
	return &responses.PermissionResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Action:      p.Action,
		RoleLevel:   p.RoleLevel.String(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
