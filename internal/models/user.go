package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash *string
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
