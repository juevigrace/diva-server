package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/storage"
	pg "github.com/juevigrace/diva-server/storage/postgres/db"
)

func pgToUUID(v pgtype.UUID) uuid.UUID {
	if !v.Valid {
		return uuid.Nil
	}
	return v.Bytes
}

func pgToUUIDPtr(v pgtype.UUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	u := uuid.UUID(v.Bytes)
	return &u
}

func pgToTime(v pgtype.Timestamptz) int64 {
	if !v.Valid {
		return 0
	}
	return v.Time.UnixMilli()
}

func pgToTimePtr(v pgtype.Timestamptz) *int64 {
	if !v.Valid {
		return nil
	}
	ms := v.Time.UnixMilli()
	return &ms
}

func DivaUserToStorage(src *pg.DivaUser) *storage.DivaUser {
	return &storage.DivaUser{
		ID:           pgToUUID(src.ID),
		Username:     src.Username,
		Email:        src.Email,
		PhoneNumber:  src.PhoneNumber,
		PasswordHash: src.PasswordHash,
		Role:         storage.RoleType(src.Role),
		CreatedAt:    pgToTime(src.CreatedAt),
		UpdatedAt:    pgToTime(src.UpdatedAt),
		DeletedAt:    pgToTimePtr(src.DeletedAt),
	}
}

func DivaPermissionToStorage(src *pg.DivaPermission) *storage.DivaPermission {
	return &storage.DivaPermission{
		ID:          pgToUUID(src.ID),
		Name:        src.Name,
		Description: src.Description,
		Action:      src.Action,
		RoleLevel:   storage.RoleType(src.RoleLevel),
		CreatedAt:   pgToTime(src.CreatedAt),
		UpdatedAt:   pgToTime(src.UpdatedAt),
		DeletedAt:   pgToTimePtr(src.DeletedAt),
	}
}

func DivaSessionToStorage(src *pg.DivaSession) *storage.DivaSession {
	return &storage.DivaSession{
		ID:           pgToUUID(src.ID),
		UserID:       pgToUUID(src.UserID),
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		Device:       src.Device,
		Type:         storage.SessionType(src.Type),
		Status:       storage.SessionStatusType(src.Status),
		IpAddress:    src.IpAddress,
		UserAgent:    src.UserAgent,
		ExpiresAt:    pgToTime(src.ExpiresAt),
		CreatedAt:    pgToTime(src.CreatedAt),
		UpdatedAt:    pgToTime(src.UpdatedAt),
	}
}

func DivaUserStateToStorage(src *pg.DivaUserState) *storage.DivaUserState {
	return &storage.DivaUserState{
		UserID:       pgToUUID(src.UserID),
		Verified:     src.Verified,
		Status:       storage.UserStatusType(src.Status),
		LastActiveAt: pgToTime(src.LastActiveAt),
		UpdatedAt:    pgToTime(src.UpdatedAt),
	}
}

func DivaUserProfileToStorage(src *pg.DivaUserProfile) *storage.DivaUserProfile {
	return &storage.DivaUserProfile{
		UserID:    pgToUUID(src.UserID),
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: pgToTimePtr(src.BirthDate),
		Alias:     src.Alias,
		Bio:       src.Bio,
		Avatar:    src.Avatar,
		UpdatedAt: pgToTime(src.UpdatedAt),
	}
}

func DivaUserPreferenceToStorage(src *pg.DivaUserPreference) *storage.DivaUserPreference {
	return &storage.DivaUserPreference{
		ID:                  pgToUUID(src.ID),
		UserID:              pgToUUID(src.UserID),
		Device:              src.Device,
		Theme:               storage.ThemeType(src.Theme),
		OnboardingCompleted: src.OnboardingCompleted,
		Language:            src.Language,
		LastSyncAt:          pgToTime(src.LastSyncAt),
		CreatedAt:           pgToTime(src.CreatedAt),
		UpdatedAt:           pgToTime(src.UpdatedAt),
	}
}

func DivaUserPermissionToStorage(src *pg.DivaUserPermission) *storage.DivaUserPermission {
	return &storage.DivaUserPermission{
		PermissionID: pgToUUID(src.PermissionID),
		UserID:       pgToUUID(src.UserID),
		GrantedBy:    pgToUUIDPtr(src.GrantedBy),
		Granted:      src.Granted,
		GrantedAt:    pgToTime(src.GrantedAt),
		ExpiresAt:    pgToTimePtr(src.ExpiresAt),
		UpdatedAt:    pgToTime(src.UpdatedAt),
	}
}

func DivaActionToStorage(src *pg.DivaAction) *storage.DivaAction {
	return &storage.DivaAction{
		ID:     pgToUUID(src.ID),
		Name:   src.Name,
		UserID: pgToUUID(src.UserID),
	}
}

func DivaActionVerificationToStorage(src *pg.GetUserVerificationRow) *storage.DivaActionVerification {
	return &storage.DivaActionVerification{
		ActionID:  pgToUUID(src.ActionID),
		Token:     src.Token,
		Verified:  src.Verified,
		ExpiresAt: pgToTime(src.ExpiresAt),
		UsedAt:    pgToTimePtr(src.UsedAt),
	}
}

func CreateUserParamsFromStorage(src *storage.CreateUserParams) *pg.CreateUserParams {
	return &pg.CreateUserParams{
		ID:           pgtype.UUID{Bytes: src.ID, Valid: true},
		Username:     src.Username,
		Email:        src.Email,
		PasswordHash: src.PasswordHash,
		Role:         pg.RoleType(src.Role),
	}
}

func CreatePermissionParamsFromStorage(src *storage.CreatePermissionParams) *pg.CreatePermissionParams {
	return &pg.CreatePermissionParams{
		ID:          pgtype.UUID{Bytes: src.ID, Valid: true},
		Name:        src.Name,
		Description: src.Description,
		Action:      src.Action,
		RoleLevel:   pg.RoleType(src.RoleLevel),
	}
}

func CreateSessionParamsFromStorage(src *storage.CreateSessionParams) *pg.CreateSessionParams {
	return &pg.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: src.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: src.UserID, Valid: true},
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		Device:       src.Device,
		Type:         pg.SessionType(src.Type),
		Status:       pg.SessionStatusType(src.Status),
		IpAddress:    src.IpAddress,
		UserAgent:    src.UserAgent,
		ExpiresAt:    pgtype.Timestamptz{Time: time.UnixMilli(src.ExpiresAt), Valid: true},
	}
}

func CreateUserPermissionParamsFromStorage(src *storage.CreateUserPermissionParams) *pg.CreateUserPermissionParams {
	var grantedBy pgtype.UUID
	if src.GrantedBy != nil {
		grantedBy = pgtype.UUID{Bytes: *src.GrantedBy, Valid: true}
	}
	var expiresAt pgtype.Timestamptz
	if src.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: time.UnixMilli(*src.ExpiresAt), Valid: true}
	}
	return &pg.CreateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: src.PermissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: src.UserID, Valid: true},
		GrantedBy:    grantedBy,
		Granted:      src.Granted,
		ExpiresAt:    expiresAt,
	}
}

func CreateUserPreferencesParamsFromStorage(src *storage.CreateUserPreferencesParams) *pg.CreateUserPreferencesParams {
	return &pg.CreateUserPreferencesParams{
		ID:                  pgtype.UUID{Bytes: src.ID, Valid: true},
		UserID:              pgtype.UUID{Bytes: src.UserID, Valid: true},
		Device:              src.Device,
		Theme:               pg.ThemeType(src.Theme),
		OnboardingCompleted: src.OnboardingCompleted,
		Language:            src.Language,
	}
}

func CreateUserProfileParamsFromStorage(src *storage.CreateUserProfileParams) *pg.CreateUserProfileParams {
	var birthDate pgtype.Timestamptz
	if src.BirthDate != nil {
		birthDate = pgtype.Timestamptz{Time: time.UnixMilli(*src.BirthDate), Valid: true}
	}
	return &pg.CreateUserProfileParams{
		UserID:    pgtype.UUID{Bytes: src.UserID, Valid: true},
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: birthDate,
		Alias:     src.Alias,
		Bio:       src.Bio,
	}
}

func CreateUserStateParamsFromStorage(src *storage.CreateUserStateParams) *pg.CreateUserStateParams {
	return &pg.CreateUserStateParams{
		UserID:   pgtype.UUID{Bytes: src.UserID, Valid: true},
		Verified: src.Verified,
		Status:   pg.UserStatusType(src.Status),
	}
}

func CreateUserActionParamsFromStorage(src *storage.CreateUserActionParams) *pg.CreateUserActionParams {
	return &pg.CreateUserActionParams{
		ID:     pgtype.UUID{Bytes: src.ID, Valid: true},
		Name:   src.Name,
		UserID: pgtype.UUID{Bytes: src.UserID, Valid: true},
	}
}

func CreateUserVerificationParamsFromStorage(src *storage.CreateUserVerificationParams) *pg.CreateUserVerificationParams {
	return &pg.CreateUserVerificationParams{
		ActionID:  pgtype.UUID{Bytes: src.ActionID, Valid: true},
		Token:     src.Token,
		ExpiresAt: pgtype.Timestamptz{Time: time.UnixMilli(src.ExpiresAt), Valid: true},
	}
}

func UpdateUserPermissionParamsFromStorage(src *storage.UpdateUserPermissionParams) *pg.UpdateUserPermissionParams {
	var expiresAt pgtype.Timestamptz
	if src.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: time.UnixMilli(*src.ExpiresAt), Valid: true}
	}
	return &pg.UpdateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: src.PermissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: src.UserID, Valid: true},
		Granted:      src.Granted,
		ExpiresAt:    expiresAt,
	}
}

func UpdateUserPreferencesParamsFromStorage(src *storage.UpdateUserPreferencesParams) *pg.UpdateUserPreferencesParams {
	return &pg.UpdateUserPreferencesParams{
		ID:       pgtype.UUID{Bytes: src.ID, Valid: true},
		Theme:    pg.ThemeType(src.Theme),
		Language: src.Language,
	}
}

func UpdateUserProfileParamsFromStorage(src *storage.UpdateUserProfileParams) *pg.UpdateUserProfileParams {
	var birthDate pgtype.Timestamptz
	if src.BirthDate != nil {
		birthDate = pgtype.Timestamptz{Time: time.UnixMilli(*src.BirthDate), Valid: true}
	}
	return &pg.UpdateUserProfileParams{
		UserID:    pgtype.UUID{Bytes: src.UserID, Valid: true},
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: birthDate,
		Alias:     src.Alias,
		Bio:       src.Bio,
	}
}

func UpdateSessionParamsFromStorage(src *storage.UpdateSessionParams) *pg.UpdateSessionParams {
	return &pg.UpdateSessionParams{
		ID:           pgtype.UUID{Bytes: src.ID, Valid: true},
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		IpAddress:    src.IpAddress,
		ExpiresAt:    pgtype.Timestamptz{Time: time.UnixMilli(src.ExpiresAt), Valid: true},
	}
}

func UpdatePermissionParamsFromStorage(src *storage.UpdatePermissionParams) *pg.UpdatePermissionParams {
	return &pg.UpdatePermissionParams{
		ID:          pgtype.UUID{Bytes: src.ID, Valid: true},
		Name:        src.Name,
		Description: src.Description,
	}
}

func ListUsersParamsFromStorage(src *storage.ListUsersParams) *pg.ListUsersParams {
	return &pg.ListUsersParams{
		Limit:  int32(src.Limit),
		Offset: int32(src.Offset),
	}
}

func ListPermissionsParamsFromStorage(src *storage.ListPermissionsParams) *pg.ListPermissionsParams {
	return &pg.ListPermissionsParams{
		Limit:  int32(src.Limit),
		Offset: int32(src.Offset),
	}
}

func UpdateEmailParamsFromStorage(src *storage.UpdateEmailParams) *pg.UpdateEmailParams {
	return &pg.UpdateEmailParams{
		Email: src.Email,
		ID:    pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdatePasswordParamsFromStorage(src *storage.UpdatePasswordParams) *pg.UpdatePasswordParams {
	return &pg.UpdatePasswordParams{
		PasswordHash: src.PasswordHash,
		ID:           pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdatePhoneNumberParamsFromStorage(src *storage.UpdatePhoneNumberParams) *pg.UpdatePhoneNumberParams {
	return &pg.UpdatePhoneNumberParams{
		PhoneNumber: src.PhoneNumber,
		ID:          pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdateRoleParamsFromStorage(src *storage.UpdateRoleParams) *pg.UpdateRoleParams {
	return &pg.UpdateRoleParams{
		Role: pg.RoleType(src.Role),
		ID:   pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdateUsernameParamsFromStorage(src *storage.UpdateUsernameParams) *pg.UpdateUsernameParams {
	return &pg.UpdateUsernameParams{
		Username: src.Username,
		ID:       pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdatePermissionActionParamsFromStorage(src *storage.UpdatePermissionActionParams) *pg.UpdatePermissionActionParams {
	return &pg.UpdatePermissionActionParams{
		Action: src.Action,
		ID:     pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdatePermissionRoleLevelParamsFromStorage(src *storage.UpdatePermissionRoleLevelParams) *pg.UpdatePermissionRoleLevelParams {
	return &pg.UpdatePermissionRoleLevelParams{
		RoleLevel: pg.RoleType(src.RoleLevel),
		ID:        pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdateSessionStatusParamsFromStorage(src *storage.UpdateSessionStatusParams) *pg.UpdateSessionStatusParams {
	return &pg.UpdateSessionStatusParams{
		Status: pg.SessionStatusType(src.Status),
		ID:     pgtype.UUID{Bytes: src.ID, Valid: true},
	}
}

func UpdateUserStatusParamsFromStorage(src *storage.UpdateUserStatusParams) *pg.UpdateUserStatusParams {
	return &pg.UpdateUserStatusParams{
		Status: pg.UserStatusType(src.Status),
		UserID: pgtype.UUID{Bytes: src.UserID, Valid: true},
	}
}

func UpdateUserVerifiedParamsFromStorage(src *storage.UpdateUserVerifiedParams) *pg.UpdateUserVerifiedParams {
	return &pg.UpdateUserVerifiedParams{
		Verified: src.Verified,
		UserID:   pgtype.UUID{Bytes: src.UserID, Valid: true},
	}
}

func UpdateUserProfileAvatarParamsFromStorage(src *storage.UpdateUserProfileAvatarParams) *pg.UpdateUserProfileAvatarParams {
	return &pg.UpdateUserProfileAvatarParams{
		Avatar: src.Avatar,
		UserID: pgtype.UUID{Bytes: src.UserID, Valid: true},
	}
}

func DeleteUserPermissionParamsFromStorage(src *storage.DeleteUserPermissionParams) *pg.DeleteUserPermissionParams {
	return &pg.DeleteUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: src.PermissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: src.UserID, Valid: true},
	}
}

func GetUserPermissionParamsFromStorage(src *storage.GetUserPermissionParams) *pg.GetUserPermissionParams {
	return &pg.GetUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: src.PermissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: src.UserID, Valid: true},
	}
}

func GetUserPermissionByNameParamsFromStorage(src *storage.GetUserPermissionByNameParams) *pg.GetUserPermissionByNameParams {
	return &pg.GetUserPermissionByNameParams{
		UserID: pgtype.UUID{Bytes: src.UserID, Valid: true},
		Name:   src.Name,
	}
}

func GetUserActionByUserAndNameParamsFromStorage(src *storage.GetUserActionByUserAndNameParams) *pg.GetUserActionByUserAndNameParams {
	return &pg.GetUserActionByUserAndNameParams{
		UserID: pgtype.UUID{Bytes: src.UserID, Valid: true},
		Name:   src.Name,
	}
}

func UpdateUserVerificationParamsFromStorage(src *storage.UpdateUserVerificationParams) *pg.UpdateUserVerificationParams {
	return &pg.UpdateUserVerificationParams{
		Verified: src.Verified,
		ActionID: pgtype.UUID{Bytes: src.ActionID, Valid: true},
	}
}
