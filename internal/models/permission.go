package models

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	Resource    string
	Action      string
	RoleLevel   Role
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   *int64
}
