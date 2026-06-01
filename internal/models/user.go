package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage/db"
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

func (up *UserPermission) DBCreate(userID uuid.UUID) *db.CreateUserPermissionParams {
	return &db.CreateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: up.Permission.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		GrantedBy:    ToUUIDPtr(up.GrantedBy),
		Granted:      up.Granted,
		GrantedAt:    ToTimestamptzPtr(up.GrantedAt),
		ExpiresAt:    ToTimestamptzPtr(up.ExpiresAt),
	}
}

func (up *UserPreferences) DBCreate(id uuid.UUID) *db.CreateUserPreferencesParams {
	return &db.CreateUserPreferencesParams{
		ID:                  pgtype.UUID{Bytes: up.ID, Valid: true},
		UserID:              pgtype.UUID{Bytes: id, Valid: true},
		Device:              up.Device,
		Theme:               up.Theme.ToDB(),
		OnboardingCompleted: up.OnboardingCompleted,
		Language:            up.Language,
	}
}

func (up *UserPreferences) DBUpdate() *db.UpdateUserPreferencesParams {
	return &db.UpdateUserPreferencesParams{
		ID:       pgtype.UUID{Bytes: up.ID, Valid: true},
		Theme:    up.Theme.ToDB(),
		Language: up.Language,
	}
}

func (up *UserProfile) DBCreate(id uuid.UUID) *db.CreateUserProfileParams {
	return &db.CreateUserProfileParams{
		UserID:    pgtype.UUID{Bytes: id, Valid: true},
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: pgtype.Timestamptz{Time: time.UnixMilli(up.BirthDate), Valid: true},
		Alias:     up.Alias,
		Bio:       up.Bio,
	}
}

func (up *UserProfile) DBUpdate(id uuid.UUID) *db.UpdateUserProfileParams {
	return &db.UpdateUserProfileParams{
		UserID:    pgtype.UUID{Bytes: id, Valid: true},
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: pgtype.Timestamptz{Time: time.UnixMilli(up.BirthDate), Valid: true},
		Alias:     up.Alias,
		Bio:       up.Bio,
	}
}

func UserPermissionFromDB(row *db.DivaUserPermission, perm *Permission) *UserPermission {
	return &UserPermission{
		Permission: *perm,
		GrantedBy:  FromPGUUIDPtr(row.GrantedBy),
		Granted:    row.Granted,
		GrantedAt:  ToInt64Ptr(row.GrantedAt),
		ExpiresAt:  ToInt64Ptr(row.ExpiresAt),
		UpdatedAt:  row.GrantedAt.Time.UnixMilli(),
	}
}

func UserPrefsFromDB(row *db.DivaUserPreference) *UserPreferences {
	return &UserPreferences{
		ID:                  row.ID.Bytes,
		Device:              row.Device,
		Theme:               ThemeFromDB(row.Theme),
		OnboardingCompleted: row.OnboardingCompleted,
		Language:            row.Language,
		LastSyncAt:          row.LastSyncAt.Time.UnixMilli(),
		CreatedAt:           row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:           row.UpdatedAt.Time.UnixMilli(),
	}
}

func UserProfileFromDB(row *db.DivaUserProfile) *UserProfile {
	return &UserProfile{
		FirstName: row.FirstName,
		LastName:  row.LastName,
		BirthDate: row.BirthDate.Time.UnixMilli(),
		Alias:     row.Alias,
		Bio:       row.Bio,
		Avatar:    row.Avatar,
	}
}
