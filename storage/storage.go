package storage

import "context"

type Storage interface {
	Close() error
	Health(ctx context.Context) HealthResult
	UserStore() UserStore
	PermissionStore() PermissionStore
	SessionStore() SessionStore
	UserStateStore() UserStateStore
	UserProfileStore() UserProfileStore
	UserPreferenceStore() UserPreferenceStore
	UserPermissionStore() UserPermissionStore
	UserActionStore() UserActionStore
	UserVerificationStore() UserVerificationStore
}
