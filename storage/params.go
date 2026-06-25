package storage

import "github.com/google/uuid"

type CreateUserParams struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	Role         RoleType
}

type CreatePermissionParams struct {
	ID          uuid.UUID
	Name        string
	Description string
	Action      string
	RoleLevel   RoleType
}

type CreateSessionParams struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	AccessToken     string
	RefreshToken    string
	Device          string
	Type            SessionType
	Status          SessionStatusType
	IpAddress       string
	UserAgent       string
	AccessExpiresAt int64
	RefreshExpiresAt int64
}

type CreateUserPermissionParams struct {
	PermissionID uuid.UUID
	UserID       uuid.UUID
	GrantedBy    *uuid.UUID
	Granted      bool
	ExpiresAt    *int64
}

type CreateUserPreferencesParams struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Device              string
	Theme               ThemeType
	OnboardingCompleted bool
	Language            string
}

type CreateUserProfileParams struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	BirthDate *int64
	Alias     string
	Bio       string
}

type CreateUserStateParams struct {
	UserID   uuid.UUID
	Verified bool
	Status   UserStatusType
}

type CreateUserActionParams struct {
	ID     uuid.UUID
	Name   string
	UserID uuid.UUID
}

type CreateUserVerificationParams struct {
	ActionID  uuid.UUID
	Token     string
	ExpiresAt int64
}

type UpdateUserPermissionParams struct {
	PermissionID uuid.UUID
	UserID       uuid.UUID
	Granted      bool
	ExpiresAt    *int64
}

type UpdateUserPreferencesParams struct {
	ID       uuid.UUID
	Theme    ThemeType
	Language string
}

type UpdateUserProfileParams struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	BirthDate *int64
	Alias     string
	Bio       string
}

type UpdateSessionParams struct {
	ID              uuid.UUID
	AccessToken     string
	RefreshToken    string
	IpAddress       string
	AccessExpiresAt int64
	RefreshExpiresAt int64
}

type UpdatePermissionParams struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type ListUsersParams struct {
	Limit  int64
	Offset int64
}

type ListPermissionsParams struct {
	Limit  int64
	Offset int64
}

type UpdateEmailParams struct {
	Email string
	ID    uuid.UUID
}

type UpdatePasswordParams struct {
	PasswordHash string
	ID           uuid.UUID
}

type UpdatePhoneNumberParams struct {
	PhoneNumber string
	ID          uuid.UUID
}

type UpdateRoleParams struct {
	Role RoleType
	ID   uuid.UUID
}

type UpdateUsernameParams struct {
	Username string
	ID       uuid.UUID
}

type UpdatePermissionActionParams struct {
	Action string
	ID     uuid.UUID
}

type UpdatePermissionRoleLevelParams struct {
	RoleLevel RoleType
	ID        uuid.UUID
}

type UpdateSessionStatusParams struct {
	Status SessionStatusType
	ID     uuid.UUID
}

type UpdateUserStatusParams struct {
	Status UserStatusType
	UserID uuid.UUID
}

type UpdateUserVerifiedParams struct {
	Verified bool
	UserID   uuid.UUID
}

type UpdateUserProfileAvatarParams struct {
	Avatar string
	UserID uuid.UUID
}

type DeleteUserPermissionParams struct {
	PermissionID uuid.UUID
	UserID       uuid.UUID
}

type GetUserPermissionParams struct {
	PermissionID uuid.UUID
	UserID       uuid.UUID
}

type GetUserPermissionByNameParams struct {
	UserID uuid.UUID
	Name   string
}

type GetUserActionByUserAndNameParams struct {
	UserID uuid.UUID
	Name   string
}

type UpdateUserVerificationParams struct {
	Verified bool
	ActionID uuid.UUID
}
