package storage

import (
	"context"

	"github.com/google/uuid"
)

type UserStore interface {
	CreateUser(ctx context.Context, arg *CreateUserParams) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*DivaUser, error)
	GetUserByEmail(ctx context.Context, email string) (*DivaUser, error)
	GetUserByUsername(ctx context.Context, username string) (*DivaUser, error)
	GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*DivaUser, error)
	ListUsers(ctx context.Context, arg *ListUsersParams) ([]DivaUser, error)
	CountUsers(ctx context.Context) (int64, error)
	UpdateUsername(ctx context.Context, arg *UpdateUsernameParams) error
	UpdateEmail(ctx context.Context, arg *UpdateEmailParams) error
	UpdatePassword(ctx context.Context, arg *UpdatePasswordParams) error
	UpdatePhoneNumber(ctx context.Context, arg *UpdatePhoneNumberParams) error
	UpdateRole(ctx context.Context, arg *UpdateRoleParams) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SoftDeleteUser(ctx context.Context, id uuid.UUID) error
	RestoreUser(ctx context.Context, id uuid.UUID) error
}

type PermissionStore interface {
	CreatePermission(ctx context.Context, arg *CreatePermissionParams) error
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*DivaPermission, error)
	GetPermissionByName(ctx context.Context, action string) (*DivaPermission, error)
	ListPermissions(ctx context.Context, arg *ListPermissionsParams) ([]DivaPermission, error)
	CountPermissions(ctx context.Context) (int64, error)
	UpdatePermission(ctx context.Context, arg *UpdatePermissionParams) error
	UpdatePermissionAction(ctx context.Context, arg *UpdatePermissionActionParams) error
	UpdatePermissionRoleLevel(ctx context.Context, arg *UpdatePermissionRoleLevelParams) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	SoftDeletePermission(ctx context.Context, id uuid.UUID) error
	RestorePermission(ctx context.Context, id uuid.UUID) error
}

type SessionStore interface {
	CreateSession(ctx context.Context, arg *CreateSessionParams) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*DivaSession, error)
	ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]DivaSession, error)
	UpdateSession(ctx context.Context, arg *UpdateSessionParams) error
	UpdateSessionStatus(ctx context.Context, arg *UpdateSessionStatusParams) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteSessionsByUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
}

type UserStateStore interface {
	CreateUserState(ctx context.Context, arg *CreateUserStateParams) error
	GetUserStateByUserID(ctx context.Context, userID uuid.UUID) (*DivaUserState, error)
	UpdateLastActiveAt(ctx context.Context, userID uuid.UUID) error
	UpdateUserStatus(ctx context.Context, arg *UpdateUserStatusParams) error
	UpdateUserVerified(ctx context.Context, arg *UpdateUserVerifiedParams) error
}

type UserProfileStore interface {
	CreateUserProfile(ctx context.Context, arg *CreateUserProfileParams) error
	GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*DivaUserProfile, error)
	UpdateUserProfile(ctx context.Context, arg *UpdateUserProfileParams) error
	UpdateUserProfileAvatar(ctx context.Context, arg *UpdateUserProfileAvatarParams) error
}

type UserPreferenceStore interface {
	CreateUserPreferences(ctx context.Context, arg *CreateUserPreferencesParams) error
	GetPreferencesByID(ctx context.Context, id uuid.UUID) (*DivaUserPreference, error)
	GetPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]DivaUserPreference, error)
	UpdateUserPreferences(ctx context.Context, arg *UpdateUserPreferencesParams) error
}

type UserPermissionStore interface {
	CreateUserPermission(ctx context.Context, arg *CreateUserPermissionParams) error
	GetUserPermission(ctx context.Context, arg *GetUserPermissionParams) (*DivaUserPermission, error)
	GetUserPermissionByName(ctx context.Context, arg *GetUserPermissionByNameParams) (*DivaUserPermission, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]DivaUserPermission, error)
	UpdateUserPermission(ctx context.Context, arg *UpdateUserPermissionParams) error
	DeleteUserPermission(ctx context.Context, arg *DeleteUserPermissionParams) error
}

type UserActionStore interface {
	CreateUserAction(ctx context.Context, arg *CreateUserActionParams) error
	GetUserActionByID(ctx context.Context, id uuid.UUID) (*DivaAction, error)
	GetUserActionByUserAndName(ctx context.Context, arg *GetUserActionByUserAndNameParams) (*DivaAction, error)
	ListActionsByUser(ctx context.Context, userID uuid.UUID) ([]DivaAction, error)
	DeleteUserAction(ctx context.Context, id uuid.UUID) error
	DeleteUserActionByUser(ctx context.Context, userID uuid.UUID) error
}

type UserVerificationStore interface {
	CreateUserVerification(ctx context.Context, arg *CreateUserVerificationParams) error
	GetUserVerification(ctx context.Context, actionID uuid.UUID) (*DivaActionVerification, error)
	UpdateUserVerification(ctx context.Context, arg *UpdateUserVerificationParams) error
	DeleteUserVerification(ctx context.Context, actionID uuid.UUID) error
}
