package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSamePassword       = errors.New("passwords are the same")
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PhoneNumber  string
	PasswordHash string
	Verified     bool
	Role         Role
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    *int64
	Profile      *UserProfile
	Actions      []UserAction
	Permissions  map[PermissionAction]UserPermission
	Preferences  UserPreferences
}

type UserProfile struct {
	FirstName string
	LastName  string
	BirthDate int64
	Alias     string
	Avatar    string
	Bio       string
}

type UserPermission struct {
	Permission Permission
	GrantedBy  *uuid.UUID
	Granted    bool
	GrantedAt  *int64
	// TODO: change expiration time for enum with fixed times
	ExpiresAt *int64
	UpdatedAt int64
}

type UserPreferences struct {
	ID                  uuid.UUID
	Device              string
	Theme               Theme
	OnboardingCompleted bool
	Language            string
	LastSyncAt          int64
	CreatedAt           int64
	UpdatedAt           int64
}

func (u *User) Response() *responses.UserResponse {
	return &responses.UserResponse{
		ID:          u.ID.String(),
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Verified:    u.Verified,
		Role:        u.Role.String(),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
}

func (up *UserProfile) Response(id *uuid.UUID) *responses.UserProfileResponse {
	var userID *string = new(string)
	if id != nil {
		*userID = id.String()
	}
	return &responses.UserProfileResponse{
		UserID:    userID,
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: up.BirthDate,
		Alias:     up.Alias,
		Avatar:    up.Avatar,
		Bio:       up.Bio,
	}
}

func (up *UserPermission) Response(id *uuid.UUID) *responses.UserPermissionResponse {
	var userID *string = new(string)
	if id != nil {
		*userID = id.String()
	}

	var grantedBy *string = new(string)
	if up.GrantedBy != nil {
		*grantedBy = up.GrantedBy.String()
	}

	return &responses.UserPermissionResponse{
		UserID:       userID,
		PermissionID: up.Permission.ID.String(),
		GrantedBy:    grantedBy,
		Granted:      up.Granted,
		GrantedAt:    up.GrantedAt,
		ExpiresAt:    up.ExpiresAt,
		UpdatedAt:    up.UpdatedAt,
	}
}

func (up *UserPreferences) Response(id *uuid.UUID) *responses.UserPreferencesResponse {
	var userID *string = new(string)
	if id != nil {
		*userID = id.String()
	}
	return &responses.UserPreferencesResponse{
		UserID:              userID,
		Id:                  up.ID.String(),
		Theme:               up.Theme.String(),
		OnboardingCompleted: up.OnboardingCompleted,
		Language:            up.Language,
		LastSyncAt:          up.LastSyncAt,
		CreatedAt:           up.CreatedAt,
		UpdatedAt:           up.UpdatedAt,
	}
}
