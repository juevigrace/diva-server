package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/storage"
	pg "github.com/juevigrace/diva-server/storage/postgres/db"
)

type UserStore struct {
	q *pg.Queries
}

func NewUserStore(q *pg.Queries) *UserStore {
	return &UserStore{q: q}
}

func (s *UserStore) CreateUser(ctx context.Context, arg *storage.CreateUserParams) error {
	params := CreateUserParamsFromStorage(arg)
	return s.q.CreateUser(ctx, *params)
}

func (s *UserStore) GetUserByID(ctx context.Context, id uuid.UUID) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(&u), nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(&u), nil
}

func (s *UserStore) GetUserByUsername(ctx context.Context, username string) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(&u), nil
}

func (s *UserStore) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByUsernameOrEmail(ctx, identifier)
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(&u), nil
}

func (s *UserStore) ListUsers(ctx context.Context, arg *storage.ListUsersParams) ([]storage.DivaUser, error) {
	params := ListUsersParamsFromStorage(arg)
	rows, err := s.q.ListUsers(ctx, *params)
	if err != nil {
		return nil, err
	}
	users := make([]storage.DivaUser, len(rows))
	for i := range rows {
		users[i] = *DivaUserToStorage(&rows[i])
	}
	return users, nil
}

func (s *UserStore) CountUsers(ctx context.Context) (int64, error) {
	return s.q.CountUsers(ctx)
}

func (s *UserStore) UpdateUsername(ctx context.Context, arg *storage.UpdateUsernameParams) error {
	params := UpdateUsernameParamsFromStorage(arg)
	return s.q.UpdateUsername(ctx, *params)
}

func (s *UserStore) UpdateEmail(ctx context.Context, arg *storage.UpdateEmailParams) error {
	params := UpdateEmailParamsFromStorage(arg)
	return s.q.UpdateEmail(ctx, *params)
}

func (s *UserStore) UpdatePassword(ctx context.Context, arg *storage.UpdatePasswordParams) error {
	params := UpdatePasswordParamsFromStorage(arg)
	return s.q.UpdatePassword(ctx, *params)
}

func (s *UserStore) UpdatePhoneNumber(ctx context.Context, arg *storage.UpdatePhoneNumberParams) error {
	params := UpdatePhoneNumberParamsFromStorage(arg)
	return s.q.UpdatePhoneNumber(ctx, *params)
}

func (s *UserStore) UpdateRole(ctx context.Context, arg *storage.UpdateRoleParams) error {
	params := UpdateRoleParamsFromStorage(arg)
	return s.q.UpdateRole(ctx, *params)
}

func (s *UserStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (s *UserStore) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (s *UserStore) RestoreUser(ctx context.Context, id uuid.UUID) error {
	return s.q.RestoreUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

type PermissionStore struct {
	q *pg.Queries
}

func NewPermissionStore(q *pg.Queries) *PermissionStore {
	return &PermissionStore{q: q}
}

func (s *PermissionStore) CreatePermission(ctx context.Context, arg *storage.CreatePermissionParams) error {
	params := CreatePermissionParamsFromStorage(arg)
	return s.q.CreatePermission(ctx, *params)
}

func (s *PermissionStore) GetPermissionByID(ctx context.Context, id uuid.UUID) (*storage.DivaPermission, error) {
	p, err := s.q.GetPermissionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaPermissionToStorage(&p), nil
}

func (s *PermissionStore) GetPermissionByName(ctx context.Context, action string) (*storage.DivaPermission, error) {
	p, err := s.q.GetPermissionByName(ctx, action)
	if err != nil {
		return nil, err
	}
	return DivaPermissionToStorage(&p), nil
}

func (s *PermissionStore) ListPermissions(ctx context.Context, arg *storage.ListPermissionsParams) ([]storage.DivaPermission, error) {
	params := ListPermissionsParamsFromStorage(arg)
	rows, err := s.q.ListPermissions(ctx, *params)
	if err != nil {
		return nil, err
	}
	perms := make([]storage.DivaPermission, len(rows))
	for i := range rows {
		perms[i] = *DivaPermissionToStorage(&rows[i])
	}
	return perms, nil
}

func (s *PermissionStore) CountPermissions(ctx context.Context) (int64, error) {
	return s.q.CountPermissions(ctx)
}

func (s *PermissionStore) UpdatePermission(ctx context.Context, arg *storage.UpdatePermissionParams) error {
	params := UpdatePermissionParamsFromStorage(arg)
	return s.q.UpdatePermission(ctx, *params)
}

func (s *PermissionStore) UpdatePermissionAction(ctx context.Context, arg *storage.UpdatePermissionActionParams) error {
	params := UpdatePermissionActionParamsFromStorage(arg)
	return s.q.UpdatePermissionAction(ctx, *params)
}

func (s *PermissionStore) UpdatePermissionRoleLevel(ctx context.Context, arg *storage.UpdatePermissionRoleLevelParams) error {
	params := UpdatePermissionRoleLevelParamsFromStorage(arg)
	return s.q.UpdatePermissionRoleLevel(ctx, *params)
}

func (s *PermissionStore) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return s.q.DeletePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (s *PermissionStore) SoftDeletePermission(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeletePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (s *PermissionStore) RestorePermission(ctx context.Context, id uuid.UUID) error {
	return s.q.RestorePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

type SessionStore struct {
	q *pg.Queries
}

func NewSessionStore(q *pg.Queries) *SessionStore {
	return &SessionStore{q: q}
}

func (s *SessionStore) CreateSession(ctx context.Context, arg *storage.CreateSessionParams) error {
	params := CreateSessionParamsFromStorage(arg)
	return s.q.CreateSession(ctx, *params)
}

func (s *SessionStore) GetSessionByID(ctx context.Context, id uuid.UUID) (*storage.DivaSession, error) {
	ss, err := s.q.GetSessionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaSessionToStorage(&ss), nil
}

func (s *SessionStore) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]storage.DivaSession, error) {
	rows, err := s.q.ListSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	sessions := make([]storage.DivaSession, len(rows))
	for i := range rows {
		sessions[i] = *DivaSessionToStorage(&rows[i])
	}
	return sessions, nil
}

func (s *SessionStore) UpdateSession(ctx context.Context, arg *storage.UpdateSessionParams) error {
	params := UpdateSessionParamsFromStorage(arg)
	return s.q.UpdateSession(ctx, *params)
}

func (s *SessionStore) UpdateSessionStatus(ctx context.Context, arg *storage.UpdateSessionStatusParams) error {
	params := UpdateSessionStatusParamsFromStorage(arg)
	return s.q.UpdateSessionStatus(ctx, *params)
}

func (s *SessionStore) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteSession(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (s *SessionStore) DeleteSessionsByUser(ctx context.Context, userID uuid.UUID) error {
	return s.q.DeleteSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (s *SessionStore) DeleteExpiredSessions(ctx context.Context) error {
	return s.q.DeleteExpiredSessions(ctx)
}

type UserStateStore struct {
	q *pg.Queries
}

func NewUserStateStore(q *pg.Queries) *UserStateStore {
	return &UserStateStore{q: q}
}

func (s *UserStateStore) CreateUserState(ctx context.Context, arg *storage.CreateUserStateParams) error {
	params := CreateUserStateParamsFromStorage(arg)
	return s.q.CreateUserState(ctx, *params)
}

func (s *UserStateStore) GetUserStateByUserID(ctx context.Context, userID uuid.UUID) (*storage.DivaUserState, error) {
	us, err := s.q.GetUserStateByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaUserStateToStorage(&us), nil
}

func (s *UserStateStore) UpdateLastActiveAt(ctx context.Context, userID uuid.UUID) error {
	return s.q.UpdateLastActiveAt(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (s *UserStateStore) UpdateUserStatus(ctx context.Context, arg *storage.UpdateUserStatusParams) error {
	params := UpdateUserStatusParamsFromStorage(arg)
	return s.q.UpdateUserStatus(ctx, *params)
}

func (s *UserStateStore) UpdateUserVerified(ctx context.Context, arg *storage.UpdateUserVerifiedParams) error {
	params := UpdateUserVerifiedParamsFromStorage(arg)
	return s.q.UpdateUserVerified(ctx, *params)
}

type UserProfileStore struct {
	q *pg.Queries
}

func NewUserProfileStore(q *pg.Queries) *UserProfileStore {
	return &UserProfileStore{q: q}
}

func (s *UserProfileStore) CreateUserProfile(ctx context.Context, arg *storage.CreateUserProfileParams) error {
	params := CreateUserProfileParamsFromStorage(arg)
	return s.q.CreateUserProfile(ctx, *params)
}

func (s *UserProfileStore) GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*storage.DivaUserProfile, error) {
	p, err := s.q.GetUserProfileByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaUserProfileToStorage(&p), nil
}

func (s *UserProfileStore) UpdateUserProfile(ctx context.Context, arg *storage.UpdateUserProfileParams) error {
	params := UpdateUserProfileParamsFromStorage(arg)
	return s.q.UpdateUserProfile(ctx, *params)
}

func (s *UserProfileStore) UpdateUserProfileAvatar(ctx context.Context, arg *storage.UpdateUserProfileAvatarParams) error {
	params := UpdateUserProfileAvatarParamsFromStorage(arg)
	return s.q.UpdateUserProfileAvatar(ctx, *params)
}

type UserPreferenceStore struct {
	q *pg.Queries
}

func NewUserPreferenceStore(q *pg.Queries) *UserPreferenceStore {
	return &UserPreferenceStore{q: q}
}

func (s *UserPreferenceStore) CreateUserPreferences(ctx context.Context, arg *storage.CreateUserPreferencesParams) error {
	params := CreateUserPreferencesParamsFromStorage(arg)
	return s.q.CreateUserPreferences(ctx, *params)
}

func (s *UserPreferenceStore) GetPreferencesByID(ctx context.Context, id uuid.UUID) (*storage.DivaUserPreference, error) {
	p, err := s.q.GetPreferencesByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaUserPreferenceToStorage(&p), nil
}

func (s *UserPreferenceStore) GetPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]storage.DivaUserPreference, error) {
	rows, err := s.q.GetPreferencesByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	prefs := make([]storage.DivaUserPreference, len(rows))
	for i := range rows {
		prefs[i] = *DivaUserPreferenceToStorage(&rows[i])
	}
	return prefs, nil
}

func (s *UserPreferenceStore) UpdateUserPreferences(ctx context.Context, arg *storage.UpdateUserPreferencesParams) error {
	params := UpdateUserPreferencesParamsFromStorage(arg)
	return s.q.UpdateUserPreferences(ctx, *params)
}

type UserPermissionStore struct {
	q *pg.Queries
}

func NewUserPermissionStore(q *pg.Queries) *UserPermissionStore {
	return &UserPermissionStore{q: q}
}

func (s *UserPermissionStore) CreateUserPermission(ctx context.Context, arg *storage.CreateUserPermissionParams) error {
	params := CreateUserPermissionParamsFromStorage(arg)
	return s.q.CreateUserPermission(ctx, *params)
}

func (s *UserPermissionStore) GetUserPermission(ctx context.Context, arg *storage.GetUserPermissionParams) (*storage.DivaUserPermission, error) {
	params := GetUserPermissionParamsFromStorage(arg)
	p, err := s.q.GetUserPermission(ctx, *params)
	if err != nil {
		return nil, err
	}
	return DivaUserPermissionToStorage(&p), nil
}

func (s *UserPermissionStore) GetUserPermissionByName(ctx context.Context, arg *storage.GetUserPermissionByNameParams) (*storage.DivaUserPermission, error) {
	params := GetUserPermissionByNameParamsFromStorage(arg)
	p, err := s.q.GetUserPermissionByName(ctx, *params)
	if err != nil {
		return nil, err
	}
	return DivaUserPermissionToStorage(&p), nil
}

func (s *UserPermissionStore) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]storage.DivaUserPermission, error) {
	rows, err := s.q.GetUserPermissions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	perms := make([]storage.DivaUserPermission, len(rows))
	for i := range rows {
		perms[i] = *DivaUserPermissionToStorage(&rows[i])
	}
	return perms, nil
}

func (s *UserPermissionStore) UpdateUserPermission(ctx context.Context, arg *storage.UpdateUserPermissionParams) error {
	params := UpdateUserPermissionParamsFromStorage(arg)
	return s.q.UpdateUserPermission(ctx, *params)
}

func (s *UserPermissionStore) DeleteUserPermission(ctx context.Context, arg *storage.DeleteUserPermissionParams) error {
	params := DeleteUserPermissionParamsFromStorage(arg)
	return s.q.DeleteUserPermission(ctx, *params)
}

type UserActionStore struct {
	q *pg.Queries
}

func NewUserActionStore(q *pg.Queries) *UserActionStore {
	return &UserActionStore{q: q}
}

func (s *UserActionStore) CreateUserAction(ctx context.Context, arg *storage.CreateUserActionParams) error {
	params := CreateUserActionParamsFromStorage(arg)
	return s.q.CreateUserAction(ctx, *params)
}

func (s *UserActionStore) GetUserActionByID(ctx context.Context, id uuid.UUID) (*storage.DivaAction, error) {
	a, err := s.q.GetUserActionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaActionToStorage(&a), nil
}

func (s *UserActionStore) GetUserActionByUserAndName(ctx context.Context, arg *storage.GetUserActionByUserAndNameParams) (*storage.DivaAction, error) {
	params := GetUserActionByUserAndNameParamsFromStorage(arg)
	a, err := s.q.GetUserActionByUserAndName(ctx, *params)
	if err != nil {
		return nil, err
	}
	return DivaActionToStorage(&a), nil
}

func (s *UserActionStore) ListActionsByUser(ctx context.Context, userID uuid.UUID) ([]storage.DivaAction, error) {
	rows, err := s.q.ListActionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	actions := make([]storage.DivaAction, len(rows))
	for i := range rows {
		actions[i] = *DivaActionToStorage(&rows[i])
	}
	return actions, nil
}

func (s *UserActionStore) DeleteUserAction(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteUserAction(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (s *UserActionStore) DeleteUserActionByUser(ctx context.Context, userID uuid.UUID) error {
	return s.q.DeleteUserActionByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

type UserVerificationStore struct {
	q *pg.Queries
}

func NewUserVerificationStore(q *pg.Queries) *UserVerificationStore {
	return &UserVerificationStore{q: q}
}

func (s *UserVerificationStore) CreateUserVerification(ctx context.Context, arg *storage.CreateUserVerificationParams) error {
	params := CreateUserVerificationParamsFromStorage(arg)
	return s.q.CreateUserVerification(ctx, *params)
}

func (s *UserVerificationStore) GetUserVerification(ctx context.Context, actionID uuid.UUID) (*storage.DivaActionVerification, error) {
	v, err := s.q.GetUserVerification(ctx, pgtype.UUID{Bytes: actionID, Valid: true})
	if err != nil {
		return nil, err
	}
	return DivaActionVerificationToStorage(&v), nil
}

func (s *UserVerificationStore) UpdateUserVerification(ctx context.Context, arg *storage.UpdateUserVerificationParams) error {
	params := UpdateUserVerificationParamsFromStorage(arg)
	return s.q.UpdateUserVerification(ctx, *params)
}

func (s *UserVerificationStore) DeleteUserVerification(ctx context.Context, actionID uuid.UUID) error {
	return s.q.DeleteUserVerification(ctx, pgtype.UUID{Bytes: actionID, Valid: true})
}
