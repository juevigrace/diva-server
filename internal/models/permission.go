package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage/db"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type PermissionAction int

const (
	PERMISSION_NONE = iota
	PERMISSION_PASSWORD_UPDATE
)

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	Action      PermissionAction
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
		Action:      p.Action.String(),
		RoleLevel:   p.RoleLevel.String(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func PermissionFromDB(row *db.DivaPermission) *Permission {
	return &Permission{
		ID:          DBUUIDToUUID(row.ID),
		Name:        row.Name,
		Description: row.Description,
		Action:      PermissionActionFromString(row.Action),
		RoleLevel:   RoleFromDB(row.RoleLevel),
		CreatedAt:   DBTimeToInt(row.CreatedAt),
		UpdatedAt:   DBTimeToInt(row.UpdatedAt),
		DeletedAt:   DBTimeToIntPtr(row.DeletedAt),
	}
}

func (p PermissionAction) String() string {
	switch p {
	case PERMISSION_PASSWORD_UPDATE:
		return "PERMISSION_PASSWORD_UPDATE"
	default:
		return "PERMISSION_NONE"
	}
}

func PermissionActionFromString(s string) PermissionAction {
	switch s {
	case "PERMISSION_PASSWORD_UPDATE":
		return PERMISSION_PASSWORD_UPDATE
	default:
		return PERMISSION_NONE
	}
}
