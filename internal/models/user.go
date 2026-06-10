package models

import (
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage/db"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PhoneNumber  string
	PasswordHash string
	Verified     bool
	Role         Role
	Status       UserStatus
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    *int64
	Profile      *UserProfile
	Actions      []UserAction
	Permissions  map[PermissionAction]UserPermission
	Preferences  UserPreferences
}

type UserProfile struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	BirthDate int64
	Alias     string
	Avatar    string
	Bio       string
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
	return &responses.UserResponse{
		ID:          u.ID.String(),
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Verified:    u.Verified,
		Role:        u.Role.String(),
		Status:      u.Status.String(),
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

func (u *User) DBCreate() *db.CreateUserParams {
	return &db.CreateUserParams{
		ID:           UUIDPtrToDB(&u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role.ToDB(),
	}
}

func (up *UserPermission) DBCreate(userID uuid.UUID) *db.CreateUserPermissionParams {
	return &db.CreateUserPermissionParams{
		PermissionID: UUIDPtrToDB(&up.Permission.ID),
		UserID:       UUIDPtrToDB(&userID),
		GrantedBy:    UUIDPtrToDB(up.GrantedBy),
		Granted:      up.Granted,
		ExpiresAt:    IntPtrToDBTime(up.ExpiresAt),
	}
}

func (up *UserPermission) DBUpdate(userID uuid.UUID) *db.UpdateUserPermissionParams {
	return &db.UpdateUserPermissionParams{
		PermissionID: UUIDPtrToDB(&up.Permission.ID),
		UserID:       UUIDPtrToDB(&userID),
		Granted:      up.Granted,
		ExpiresAt:    IntPtrToDBTime(up.ExpiresAt),
	}
}

func (up *UserPreferences) DBCreate(userID uuid.UUID) *db.CreateUserPreferencesParams {
	return &db.CreateUserPreferencesParams{
		ID:                  UUIDPtrToDB(&up.ID),
		UserID:              UUIDPtrToDB(&userID),
		Device:              up.Device,
		Theme:               up.Theme.ToDB(),
		OnboardingCompleted: up.OnboardingCompleted,
		Language:            up.Language,
	}
}

func (up *UserPreferences) DBUpdate() *db.UpdateUserPreferencesParams {
	return &db.UpdateUserPreferencesParams{
		ID:       UUIDPtrToDB(&up.ID),
		Theme:    up.Theme.ToDB(),
		Language: up.Language,
	}
}

func (up *UserProfile) DBCreate(userID uuid.UUID) *db.CreateUserProfileParams {
	return &db.CreateUserProfileParams{
		UserID:    UUIDPtrToDB(&userID),
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: IntPtrToDBTime(&up.BirthDate),
		Alias:     up.Alias,
		Bio:       up.Bio,
	}
}

func (up *UserProfile) DBUpdate(userID uuid.UUID) *db.UpdateUserProfileParams {
	return &db.UpdateUserProfileParams{
		UserID:    UUIDPtrToDB(&userID),
		FirstName: up.FirstName,
		LastName:  up.LastName,
		BirthDate: IntPtrToDBTime(&up.BirthDate),
		Alias:     up.Alias,
		Bio:       up.Bio,
	}
}

func UserFromDB(row *db.DivaUser) *User {
	return &User{
		ID:           DBUUIDToUUID(row.ID),
		Username:     row.Username,
		Email:        row.Email,
		PhoneNumber:  row.PhoneNumber,
		PasswordHash: row.PasswordHash,
		Verified:     row.Verified,
		Role:         RoleFromDB(row.Role),
		Status:       UserStatusFromDB(row.Status),
		CreatedAt:    DBTimeToInt(row.CreatedAt),
		UpdatedAt:    DBTimeToInt(row.UpdatedAt),
		DeletedAt:    DBTimeToIntPtr(row.DeletedAt),
	}
}

func UserPermissionFromDB(row *db.DivaUserPermission, perm *Permission) *UserPermission {
	return &UserPermission{
		Permission: *perm,
		UserID:     DBUUIDToUUID(row.UserID),
		GrantedBy:  DBUUIDToUUIDPtr(row.GrantedBy),
		Granted:    row.Granted,
		GrantedAt:  DBTimeToInt(row.GrantedAt),
		ExpiresAt:  DBTimeToIntPtr(row.ExpiresAt),
		UpdatedAt:  DBTimeToInt(row.UpdatedAt),
	}
}

func UserPrefsFromDB(row *db.DivaUserPreference) *UserPreferences {
	return &UserPreferences{
		ID:                  row.ID.Bytes,
		UserID:              DBUUIDToUUID(row.UserID),
		Device:              row.Device,
		Theme:               ThemeFromDB(row.Theme),
		OnboardingCompleted: row.OnboardingCompleted,
		Language:            row.Language,
		LastSyncAt:          DBTimeToInt(row.LastSyncAt),
		CreatedAt:           DBTimeToInt(row.CreatedAt),
		UpdatedAt:           DBTimeToInt(row.UpdatedAt),
	}
}

func UserProfileFromDB(row *db.DivaUserProfile) *UserProfile {
	return &UserProfile{
		UserID:    DBUUIDToUUID(row.UserID),
		FirstName: row.FirstName,
		LastName:  row.LastName,
		BirthDate: DBTimeToInt(row.BirthDate),
		Alias:     row.Alias,
		Bio:       row.Bio,
		Avatar:    row.Avatar,
	}
}
