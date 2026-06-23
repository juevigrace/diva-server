package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/storage"
	sqli "github.com/juevigrace/diva-server/storage/sqlite/db"
)

func sqliteToUUID(v string) uuid.UUID {
	u, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil
	}
	return u
}

func sqliteToUUIDPtr(v string) *uuid.UUID {
	if v == "" {
		return nil
	}
	u, err := uuid.Parse(v)
	if err != nil {
		return nil
	}
	return &u
}

func sqlNullStringToUUIDPtr(v sql.NullString) *uuid.UUID {
	if !v.Valid || v.String == "" {
		return nil
	}
	u, err := uuid.Parse(v.String)
	if err != nil {
		return nil
	}
	return &u
}

func sqlNullTimeToTimePtr(v sql.NullTime) *int64 {
	if !v.Valid {
		return nil
	}
	ms := v.Time.UnixMilli()
	return &ms
}

func DivaUserToStorage(src sqli.DivaUser) *storage.DivaUser {
	return &storage.DivaUser{
		ID:           sqliteToUUID(src.ID),
		Username:     src.Username,
		Email:        src.Email,
		PhoneNumber:  src.PhoneNumber,
		PasswordHash: src.PasswordHash,
		Role:         storage.RoleType(src.Role),
		CreatedAt:    src.CreatedAt.UnixMilli(),
		UpdatedAt:    src.UpdatedAt.UnixMilli(),
		DeletedAt:    sqlNullTimeToTimePtr(src.DeletedAt),
	}
}

func DivaUserFromStorage(src storage.DivaUser) sqli.DivaUser {
	var deletedAt sql.NullTime
	if src.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: time.UnixMilli(*src.DeletedAt), Valid: true}
	}
	return sqli.DivaUser{
		ID:           src.ID.String(),
		Username:     src.Username,
		Email:        src.Email,
		PhoneNumber:  src.PhoneNumber,
		PasswordHash: src.PasswordHash,
		Role:         string(src.Role),
		CreatedAt:    time.UnixMilli(src.CreatedAt),
		UpdatedAt:    time.UnixMilli(src.UpdatedAt),
		DeletedAt:    deletedAt,
	}
}

func DivaPermissionToStorage(src sqli.DivaPermission) *storage.DivaPermission {
	return &storage.DivaPermission{
		ID:          sqliteToUUID(src.ID),
		Name:        src.Name,
		Description: src.Description,
		Action:      src.Action,
		RoleLevel:   storage.RoleType(src.RoleLevel),
		CreatedAt:   src.CreatedAt.UnixMilli(),
		UpdatedAt:   src.UpdatedAt.UnixMilli(),
		DeletedAt:   sqlNullTimeToTimePtr(src.DeletedAt),
	}
}

func DivaPermissionFromStorage(src storage.DivaPermission) sqli.DivaPermission {
	var deletedAt sql.NullTime
	if src.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: time.UnixMilli(*src.DeletedAt), Valid: true}
	}
	return sqli.DivaPermission{
		ID:          src.ID.String(),
		Name:        src.Name,
		Description: src.Description,
		Action:      src.Action,
		RoleLevel:   string(src.RoleLevel),
		CreatedAt:   time.UnixMilli(src.CreatedAt),
		UpdatedAt:   time.UnixMilli(src.UpdatedAt),
		DeletedAt:   deletedAt,
	}
}

func DivaSessionToStorage(src sqli.DivaSession) *storage.DivaSession {
	return &storage.DivaSession{
		ID:           sqliteToUUID(src.ID),
		UserID:       sqliteToUUID(src.UserID),
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		Device:       src.Device,
		Type:         storage.SessionType(src.Type),
		Status:       storage.SessionStatusType(src.Status),
		IpAddress:    src.IpAddress,
		UserAgent:    src.UserAgent,
		ExpiresAt:    src.ExpiresAt.UnixMilli(),
		CreatedAt:    src.CreatedAt.UnixMilli(),
		UpdatedAt:    src.UpdatedAt.UnixMilli(),
	}
}

func DivaSessionFromStorage(src storage.DivaSession) sqli.DivaSession {
	return sqli.DivaSession{
		ID:           src.ID.String(),
		UserID:       src.UserID.String(),
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		Device:       src.Device,
		Type:         string(src.Type),
		Status:       string(src.Status),
		IpAddress:    src.IpAddress,
		UserAgent:    src.UserAgent,
		ExpiresAt:    time.UnixMilli(src.ExpiresAt),
		CreatedAt:    time.UnixMilli(src.CreatedAt),
		UpdatedAt:    time.UnixMilli(src.UpdatedAt),
	}
}

func DivaUserStateToStorage(src sqli.DivaUserState) *storage.DivaUserState {
	return &storage.DivaUserState{
		UserID:       sqliteToUUID(src.UserID),
		Verified:     src.Verified,
		Status:       storage.UserStatusType(src.Status),
		LastActiveAt: src.LastActiveAt.UnixMilli(),
		UpdatedAt:    src.UpdatedAt.UnixMilli(),
	}
}

func DivaUserStateFromStorage(src storage.DivaUserState) sqli.DivaUserState {
	return sqli.DivaUserState{
		UserID:       src.UserID.String(),
		Verified:     src.Verified,
		Status:       string(src.Status),
		LastActiveAt: time.UnixMilli(src.LastActiveAt),
		UpdatedAt:    time.UnixMilli(src.UpdatedAt),
	}
}

func DivaUserProfileToStorage(src sqli.DivaUserProfile) *storage.DivaUserProfile {
	return &storage.DivaUserProfile{
		UserID:    sqliteToUUID(src.UserID),
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: sqlNullTimeToTimePtr(src.BirthDate),
		Alias:     src.Alias,
		Bio:       src.Bio,
		Avatar:    src.Avatar,
		UpdatedAt: src.UpdatedAt.UnixMilli(),
	}
}

func DivaUserProfileFromStorage(src storage.DivaUserProfile) sqli.DivaUserProfile {
	var birthDate sql.NullTime
	if src.BirthDate != nil {
		birthDate = sql.NullTime{Time: time.UnixMilli(*src.BirthDate), Valid: true}
	}
	return sqli.DivaUserProfile{
		UserID:    src.UserID.String(),
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: birthDate,
		Alias:     src.Alias,
		Bio:       src.Bio,
		Avatar:    src.Avatar,
		UpdatedAt: time.UnixMilli(src.UpdatedAt),
	}
}

func DivaUserPreferenceToStorage(src sqli.DivaUserPreference) *storage.DivaUserPreference {
	return &storage.DivaUserPreference{
		ID:                  sqliteToUUID(src.ID),
		UserID:              sqliteToUUID(src.UserID),
		Device:              src.Device,
		Theme:               storage.ThemeType(src.Theme),
		OnboardingCompleted: src.OnboardingCompleted,
		Language:            src.Language,
		LastSyncAt:          src.LastSyncAt.UnixMilli(),
		CreatedAt:           src.CreatedAt.UnixMilli(),
		UpdatedAt:           src.UpdatedAt.UnixMilli(),
	}
}

func DivaUserPreferenceFromStorage(src storage.DivaUserPreference) sqli.DivaUserPreference {
	return sqli.DivaUserPreference{
		ID:                  src.ID.String(),
		UserID:              src.UserID.String(),
		Device:              src.Device,
		Theme:               string(src.Theme),
		OnboardingCompleted: src.OnboardingCompleted,
		Language:            src.Language,
		LastSyncAt:          time.UnixMilli(src.LastSyncAt),
		CreatedAt:           time.UnixMilli(src.CreatedAt),
		UpdatedAt:           time.UnixMilli(src.UpdatedAt),
	}
}

func DivaUserPermissionToStorage(src sqli.DivaUserPermission) *storage.DivaUserPermission {
	return &storage.DivaUserPermission{
		PermissionID: sqliteToUUID(src.PermissionID),
		UserID:       sqliteToUUID(src.UserID),
		GrantedBy:    sqlNullStringToUUIDPtr(src.GrantedBy),
		Granted:      src.Granted,
		GrantedAt:    src.GrantedAt.UnixMilli(),
		ExpiresAt:    sqlNullTimeToTimePtr(src.ExpiresAt),
		UpdatedAt:    src.UpdatedAt.UnixMilli(),
	}
}

func DivaUserPermissionFromStorage(src storage.DivaUserPermission) sqli.DivaUserPermission {
	var grantedBy sql.NullString
	if src.GrantedBy != nil {
		grantedBy = sql.NullString{String: src.GrantedBy.String(), Valid: true}
	}
	var expiresAt sql.NullTime
	if src.ExpiresAt != nil {
		expiresAt = sql.NullTime{Time: time.UnixMilli(*src.ExpiresAt), Valid: true}
	}
	return sqli.DivaUserPermission{
		PermissionID: src.PermissionID.String(),
		UserID:       src.UserID.String(),
		GrantedBy:    grantedBy,
		Granted:      src.Granted,
		GrantedAt:    time.UnixMilli(src.GrantedAt),
		ExpiresAt:    expiresAt,
		UpdatedAt:    time.UnixMilli(src.UpdatedAt),
	}
}

func DivaActionToStorage(src sqli.DivaAction) *storage.DivaAction {
	return &storage.DivaAction{
		ID:     sqliteToUUID(src.ID),
		Name:   src.Name,
		UserID: sqliteToUUID(src.UserID),
	}
}

func DivaActionFromStorage(src storage.DivaAction) sqli.DivaAction {
	return sqli.DivaAction{
		ID:     src.ID.String(),
		Name:   src.Name,
		UserID: src.UserID.String(),
	}
}

func DivaActionVerificationToStorage(src sqli.GetUserVerificationRow) *storage.DivaActionVerification {
	return &storage.DivaActionVerification{
		ActionID:  sqliteToUUID(src.ActionID),
		Token:     src.Token,
		Verified:  src.Verified,
		ExpiresAt: src.ExpiresAt.UnixMilli(),
		UsedAt:    sqlNullTimeToTimePtr(src.UsedAt),
	}
}

func DivaActionVerificationFromStorage(src storage.DivaActionVerification) sqli.GetUserVerificationRow {
	var usedAt sql.NullTime
	if src.UsedAt != nil {
		usedAt = sql.NullTime{Time: time.UnixMilli(*src.UsedAt), Valid: true}
	}
	return sqli.GetUserVerificationRow{
		ActionID:  src.ActionID.String(),
		Token:     src.Token,
		Verified:  src.Verified,
		ExpiresAt: time.UnixMilli(src.ExpiresAt),
		UsedAt:    usedAt,
	}
}

func CreateUserParamsFromStorage(src storage.CreateUserParams) sqli.CreateUserParams {
	return sqli.CreateUserParams{
		ID:           src.ID.String(),
		Username:     src.Username,
		Email:        src.Email,
		PasswordHash: src.PasswordHash,
		Role:         string(src.Role),
	}
}

func CreatePermissionParamsFromStorage(src storage.CreatePermissionParams) sqli.CreatePermissionParams {
	return sqli.CreatePermissionParams{
		ID:          src.ID.String(),
		Name:        src.Name,
		Description: src.Description,
		Action:      src.Action,
		RoleLevel:   string(src.RoleLevel),
	}
}

func CreateSessionParamsFromStorage(src storage.CreateSessionParams) sqli.CreateSessionParams {
	return sqli.CreateSessionParams{
		ID:           src.ID.String(),
		UserID:       src.UserID.String(),
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		Device:       src.Device,
		Type:         string(src.Type),
		Status:       string(src.Status),
		IpAddress:    src.IpAddress,
		UserAgent:    src.UserAgent,
		ExpiresAt:    time.UnixMilli(src.ExpiresAt),
	}
}

func CreateUserPermissionParamsFromStorage(src storage.CreateUserPermissionParams) sqli.CreateUserPermissionParams {
	var grantedBy sql.NullString
	if src.GrantedBy != nil {
		grantedBy = sql.NullString{String: src.GrantedBy.String(), Valid: true}
	}
	var expiresAt sql.NullTime
	if src.ExpiresAt != nil {
		expiresAt = sql.NullTime{Time: time.UnixMilli(*src.ExpiresAt), Valid: true}
	}
	return sqli.CreateUserPermissionParams{
		PermissionID: src.PermissionID.String(),
		UserID:       src.UserID.String(),
		GrantedBy:    grantedBy,
		Granted:      src.Granted,
		ExpiresAt:    expiresAt,
	}
}

func CreateUserPreferencesParamsFromStorage(src storage.CreateUserPreferencesParams) sqli.CreateUserPreferencesParams {
	return sqli.CreateUserPreferencesParams{
		ID:                  src.ID.String(),
		UserID:              src.UserID.String(),
		Device:              src.Device,
		Theme:               string(src.Theme),
		OnboardingCompleted: src.OnboardingCompleted,
		Language:            src.Language,
	}
}

func CreateUserProfileParamsFromStorage(src storage.CreateUserProfileParams) sqli.CreateUserProfileParams {
	var birthDate sql.NullTime
	if src.BirthDate != nil {
		birthDate = sql.NullTime{Time: time.UnixMilli(*src.BirthDate), Valid: true}
	}
	return sqli.CreateUserProfileParams{
		UserID:    src.UserID.String(),
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: birthDate,
		Alias:     src.Alias,
		Bio:       src.Bio,
	}
}

func CreateUserStateParamsFromStorage(src storage.CreateUserStateParams) sqli.CreateUserStateParams {
	return sqli.CreateUserStateParams{
		UserID:   src.UserID.String(),
		Verified: src.Verified,
		Status:   string(src.Status),
	}
}

func CreateUserActionParamsFromStorage(src storage.CreateUserActionParams) sqli.CreateUserActionParams {
	return sqli.CreateUserActionParams{
		ID:     src.ID.String(),
		Name:   src.Name,
		UserID: src.UserID.String(),
	}
}

func CreateUserVerificationParamsFromStorage(src storage.CreateUserVerificationParams) sqli.CreateUserVerificationParams {
	return sqli.CreateUserVerificationParams{
		ActionID:  src.ActionID.String(),
		Token:     src.Token,
		ExpiresAt: time.UnixMilli(src.ExpiresAt),
	}
}

func UpdateUserPermissionParamsFromStorage(src storage.UpdateUserPermissionParams) sqli.UpdateUserPermissionParams {
	var expiresAt sql.NullTime
	if src.ExpiresAt != nil {
		expiresAt = sql.NullTime{Time: time.UnixMilli(*src.ExpiresAt), Valid: true}
	}
	return sqli.UpdateUserPermissionParams{
		PermissionID: src.PermissionID.String(),
		UserID:       src.UserID.String(),
		Granted:      src.Granted,
		ExpiresAt:    expiresAt,
	}
}

func UpdateUserPreferencesParamsFromStorage(src storage.UpdateUserPreferencesParams) sqli.UpdateUserPreferencesParams {
	return sqli.UpdateUserPreferencesParams{
		ID:       src.ID.String(),
		Theme:    string(src.Theme),
		Language: src.Language,
	}
}

func UpdateUserProfileParamsFromStorage(src storage.UpdateUserProfileParams) sqli.UpdateUserProfileParams {
	var birthDate sql.NullTime
	if src.BirthDate != nil {
		birthDate = sql.NullTime{Time: time.UnixMilli(*src.BirthDate), Valid: true}
	}
	return sqli.UpdateUserProfileParams{
		UserID:    src.UserID.String(),
		FirstName: src.FirstName,
		LastName:  src.LastName,
		BirthDate: birthDate,
		Alias:     src.Alias,
		Bio:       src.Bio,
	}
}

func UpdateSessionParamsFromStorage(src storage.UpdateSessionParams) sqli.UpdateSessionParams {
	return sqli.UpdateSessionParams{
		ID:           src.ID.String(),
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
		IpAddress:    src.IpAddress,
		ExpiresAt:    time.UnixMilli(src.ExpiresAt),
	}
}

func UpdatePermissionParamsFromStorage(src storage.UpdatePermissionParams) sqli.UpdatePermissionParams {
	return sqli.UpdatePermissionParams{
		ID:          src.ID.String(),
		Name:        src.Name,
		Description: src.Description,
	}
}

func ListUsersParamsFromStorage(src storage.ListUsersParams) sqli.ListUsersParams {
	return sqli.ListUsersParams{
		Limit:  src.Limit,
		Offset: src.Offset,
	}
}

func ListPermissionsParamsFromStorage(src storage.ListPermissionsParams) sqli.ListPermissionsParams {
	return sqli.ListPermissionsParams{
		Limit:  src.Limit,
		Offset: src.Offset,
	}
}

func UpdateEmailParamsFromStorage(src storage.UpdateEmailParams) sqli.UpdateEmailParams {
	return sqli.UpdateEmailParams{
		Email: src.Email,
		ID:    src.ID.String(),
	}
}

func UpdatePasswordParamsFromStorage(src storage.UpdatePasswordParams) sqli.UpdatePasswordParams {
	return sqli.UpdatePasswordParams{
		PasswordHash: src.PasswordHash,
		ID:           src.ID.String(),
	}
}

func UpdatePhoneNumberParamsFromStorage(src storage.UpdatePhoneNumberParams) sqli.UpdatePhoneNumberParams {
	return sqli.UpdatePhoneNumberParams{
		PhoneNumber: src.PhoneNumber,
		ID:          src.ID.String(),
	}
}

func UpdateRoleParamsFromStorage(src storage.UpdateRoleParams) sqli.UpdateRoleParams {
	return sqli.UpdateRoleParams{
		Role: string(src.Role),
		ID:   src.ID.String(),
	}
}

func UpdateUsernameParamsFromStorage(src storage.UpdateUsernameParams) sqli.UpdateUsernameParams {
	return sqli.UpdateUsernameParams{
		Username: src.Username,
		ID:       src.ID.String(),
	}
}

func UpdatePermissionActionParamsFromStorage(src storage.UpdatePermissionActionParams) sqli.UpdatePermissionActionParams {
	return sqli.UpdatePermissionActionParams{
		Action: src.Action,
		ID:     src.ID.String(),
	}
}

func UpdatePermissionRoleLevelParamsFromStorage(src storage.UpdatePermissionRoleLevelParams) sqli.UpdatePermissionRoleLevelParams {
	return sqli.UpdatePermissionRoleLevelParams{
		RoleLevel: string(src.RoleLevel),
		ID:        src.ID.String(),
	}
}

func UpdateSessionStatusParamsFromStorage(src storage.UpdateSessionStatusParams) sqli.UpdateSessionStatusParams {
	return sqli.UpdateSessionStatusParams{
		Status: string(src.Status),
		ID:     src.ID.String(),
	}
}

func UpdateUserStatusParamsFromStorage(src storage.UpdateUserStatusParams) sqli.UpdateUserStatusParams {
	return sqli.UpdateUserStatusParams{
		Status: string(src.Status),
		UserID: src.UserID.String(),
	}
}

func UpdateUserVerifiedParamsFromStorage(src storage.UpdateUserVerifiedParams) sqli.UpdateUserVerifiedParams {
	return sqli.UpdateUserVerifiedParams{
		Verified: src.Verified,
		UserID:   src.UserID.String(),
	}
}

func UpdateUserProfileAvatarParamsFromStorage(src storage.UpdateUserProfileAvatarParams) sqli.UpdateUserProfileAvatarParams {
	return sqli.UpdateUserProfileAvatarParams{
		Avatar: src.Avatar,
		UserID: src.UserID.String(),
	}
}

func DeleteUserPermissionParamsFromStorage(src storage.DeleteUserPermissionParams) sqli.DeleteUserPermissionParams {
	return sqli.DeleteUserPermissionParams{
		PermissionID: src.PermissionID.String(),
		UserID:       src.UserID.String(),
	}
}

func GetUserPermissionParamsFromStorage(src storage.GetUserPermissionParams) sqli.GetUserPermissionParams {
	return sqli.GetUserPermissionParams{
		PermissionID: src.PermissionID.String(),
		UserID:       src.UserID.String(),
	}
}

func GetUserPermissionByNameParamsFromStorage(src storage.GetUserPermissionByNameParams) sqli.GetUserPermissionByNameParams {
	return sqli.GetUserPermissionByNameParams{
		UserID: src.UserID.String(),
		Name:   src.Name,
	}
}

func GetUserActionByUserAndNameParamsFromStorage(src storage.GetUserActionByUserAndNameParams) sqli.GetUserActionByUserAndNameParams {
	return sqli.GetUserActionByUserAndNameParams{
		UserID: src.UserID.String(),
		Name:   src.Name,
	}
}

func UpdateUserVerificationParamsFromStorage(src storage.UpdateUserVerificationParams) sqli.UpdateUserVerificationParams {
	return sqli.UpdateUserVerificationParams{
		Verified: src.Verified,
		ActionID: src.ActionID.String(),
	}
}
