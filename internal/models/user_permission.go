package models

import "github.com/google/uuid"

type UserPermission struct {
	Permission uuid.UUID
	UserID     uuid.UUID
	GrantedBy  *uuid.UUID
	Granted    bool
	GrantedAt  int64
	ExpiresAt  *int64
	UpdatedAt  int64
}
