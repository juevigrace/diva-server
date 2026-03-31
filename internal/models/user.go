package models

import (
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	BirthDate    int64
	PhoneNumber  string
	Alias        string
	Avatar       string
	Bio          string
	UserVerified bool
	Role         Role
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    *int64
	Permissions  []*UserPermission
}

func (u *User) ToUserResponse() *responses.UserResponse {
	return &responses.UserResponse{
		ID:           u.ID.String(),
		Email:        u.Email,
		Username:     u.Username,
		BirthDate:    u.BirthDate,
		PhoneNumber:  u.PhoneNumber,
		Alias:        u.Alias,
		Avatar:       u.Avatar,
		Bio:          u.Bio,
		UserVerified: u.UserVerified,
		Role:         u.Role.String(),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
	}
}
