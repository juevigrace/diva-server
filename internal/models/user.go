package models

import (
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PhoneNumber  string
	PasswordHash string
	Role         Role
	State        *UserState
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    *int64
	Profile      *UserProfile
	Actions      []UserAction
	Permissions  map[PermissionAction]UserPermission
	Preferences  UserPreferences
}

type UserState struct {
	Verified     bool
	Status       UserStatus
	LastActiveAt int64
	UpdatedAt    int64
}

type UserProfile struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	BirthDate int64
	Alias     string
	Avatar    string
	Bio       string
	UpdatedAt int64
}

type UserPermission struct {
	Permission Permission
	UserID     uuid.UUID
	GrantedBy  *uuid.UUID
	Granted    bool
	GrantedAt  int64
	// TODO: change expiration time for enum with fixed times
	ExpiresAt *int64
	UpdatedAt int64
}

type UserPreferences struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Device              string
	Theme               Theme
	OnboardingCompleted bool
	Language            string
	LastSyncAt          int64
	CreatedAt           int64
	UpdatedAt           int64
}

func (u *User) Response() *responses.UserResponse {
	var state *responses.UserStateResponse
	if u.State != nil {
		state = &responses.UserStateResponse{
			Verified:     u.State.Verified,
			Status:       u.State.Status.String(),
			LastActiveAt: u.State.LastActiveAt,
		}
	}

	return &responses.UserResponse{
		ID:          u.ID.String(),
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Role:        u.Role.String(),
		State:       state,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
}

func (up *UserProfile) Response() *responses.UserProfileResponse {
	return &responses.UserProfileResponse{
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: up.BirthDate,
		Alias:     up.Alias,
		Avatar:    up.Avatar,
		Bio:       up.Bio,
	}
}

func (up *UserPermission) Response() *responses.UserPermissionResponse {
	var grantedBy *string
	if up.GrantedBy != nil {
		grantedBy = new(string)
		*grantedBy = up.GrantedBy.String()
	}

	return &responses.UserPermissionResponse{
		PermissionID: up.Permission.ID.String(),
		GrantedBy:    grantedBy,
		Granted:      up.Granted,
		GrantedAt:    up.GrantedAt,
		ExpiresAt:    up.ExpiresAt,
		UpdatedAt:    up.UpdatedAt,
	}
}

func (up *UserPreferences) Response() *responses.UserPreferencesResponse {
	return &responses.UserPreferencesResponse{
		Id:                  up.ID.String(),
		Theme:               up.Theme.String(),
		OnboardingCompleted: up.OnboardingCompleted,
		Language:            up.Language,
		LastSyncAt:          up.LastSyncAt,
		CreatedAt:           up.CreatedAt,
		UpdatedAt:           up.UpdatedAt,
	}
}

func (u *User) DBCreate() *storage.CreateUserParams {
	return &storage.CreateUserParams{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role.ToDB(),
	}
}

func (up *UserPermission) DBCreate(userID uuid.UUID) *storage.CreateUserPermissionParams {
	return &storage.CreateUserPermissionParams{
		PermissionID: up.Permission.ID,
		UserID:       userID,
		GrantedBy:    up.GrantedBy,
		Granted:      up.Granted,
		ExpiresAt:    up.ExpiresAt,
	}
}

func (up *UserPermission) DBUpdate(userID uuid.UUID) *storage.UpdateUserPermissionParams {
	return &storage.UpdateUserPermissionParams{
		PermissionID: up.Permission.ID,
		UserID:       userID,
		Granted:      up.Granted,
		ExpiresAt:    up.ExpiresAt,
	}
}

func (up *UserPreferences) DBCreate(userID uuid.UUID) *storage.CreateUserPreferencesParams {
	return &storage.CreateUserPreferencesParams{
		ID:                  up.ID,
		UserID:              userID,
		Device:              up.Device,
		Theme:               up.Theme.ToDB(),
		OnboardingCompleted: up.OnboardingCompleted,
		Language:            up.Language,
	}
}

func (up *UserPreferences) DBUpdate() *storage.UpdateUserPreferencesParams {
	return &storage.UpdateUserPreferencesParams{
		ID:       up.ID,
		Theme:    up.Theme.ToDB(),
		Language: up.Language,
	}
}

func (up *UserProfile) DBCreate(userID uuid.UUID) *storage.CreateUserProfileParams {
	return &storage.CreateUserProfileParams{
		UserID:    userID,
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: &up.BirthDate,
		Alias:     up.Alias,
		Bio:       up.Bio,
	}
}

func (up *UserProfile) DBUpdate(userID uuid.UUID) *storage.UpdateUserProfileParams {
	return &storage.UpdateUserProfileParams{
		UserID:    userID,
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: &up.BirthDate,
		Alias:     up.Alias,
		Bio:       up.Bio,
	}
}

func UserFromDB(row *storage.DivaUser) *User {
	return &User{
		ID:           row.ID,
		Username:     row.Username,
		Email:        row.Email,
		PhoneNumber:  row.PhoneNumber,
		PasswordHash: row.PasswordHash,
		Role:         RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		DeletedAt:    row.DeletedAt,
	}
}

func UserPermissionFromDB(row *storage.DivaUserPermission, perm *Permission) *UserPermission {
	return &UserPermission{
		Permission: *perm,
		UserID:     row.UserID,
		GrantedBy:  row.GrantedBy,
		Granted:    row.Granted,
		GrantedAt:  row.GrantedAt,
		ExpiresAt:  row.ExpiresAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func UserPrefsFromDB(row *storage.DivaUserPreference) *UserPreferences {
	return &UserPreferences{
		ID:                  row.ID,
		UserID:              row.UserID,
		Device:              row.Device,
		Theme:               ThemeFromDB(row.Theme),
		OnboardingCompleted: row.OnboardingCompleted,
		Language:            row.Language,
		LastSyncAt:          row.LastSyncAt,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}

func UserProfileFromDB(row *storage.DivaUserProfile) *UserProfile {
	birthDate := int64(0)
	if row.BirthDate != nil {
		birthDate = *row.BirthDate
	}
	return &UserProfile{
		UserID:    row.UserID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		BirthDate: birthDate,
		Alias:     row.Alias,
		Bio:       row.Bio,
		Avatar:    row.Avatar,
		UpdatedAt: row.UpdatedAt,
	}
}

func UserStateFromDB(row *storage.DivaUserState) *UserState {
	return &UserState{
		Verified:     row.Verified,
		Status:       UserStatusFromDB(row.Status),
		LastActiveAt: row.LastActiveAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func (us *UserState) DBCreate(userID uuid.UUID) *storage.CreateUserStateParams {
	return &storage.CreateUserStateParams{
		UserID:   userID,
		Verified: us.Verified,
		Status:   us.Status.ToDB(),
	}
}
