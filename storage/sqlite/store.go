package sqlite

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/storage"
	sqli "github.com/juevigrace/diva-server/storage/sqlite/db"
)

type UserStore struct {
	q *sqli.Queries
}

func NewUserStore(q *sqli.Queries) *UserStore {
	return &UserStore{q: q}
}

func (s *UserStore) CreateUser(ctx context.Context, arg storage.CreateUserParams) error {
	return s.q.CreateUser(ctx, CreateUserParamsFromStorage(arg))
}

func (s *UserStore) GetUserByID(ctx context.Context, id uuid.UUID) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(u), nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(u), nil
}

func (s *UserStore) GetUserByUsername(ctx context.Context, username string) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(u), nil
}

func (s *UserStore) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*storage.DivaUser, error) {
	u, err := s.q.GetUserByUsernameOrEmail(ctx, sqli.GetUserByUsernameOrEmailParams{Email: identifier, Username: identifier})
	if err != nil {
		return nil, err
	}
	return DivaUserToStorage(u), nil
}

func (s *UserStore) ListUsers(ctx context.Context, arg storage.ListUsersParams) ([]storage.DivaUser, error) {
	rows, err := s.q.ListUsers(ctx, ListUsersParamsFromStorage(arg))
	if err != nil {
		return nil, err
	}
	users := make([]storage.DivaUser, len(rows))
	for i := range rows {
		users[i] = *DivaUserToStorage(rows[i])
	}
	return users, nil
}

func (s *UserStore) CountUsers(ctx context.Context) (int64, error) {
	return s.q.CountUsers(ctx)
}

func (s *UserStore) UpdateUsername(ctx context.Context, arg storage.UpdateUsernameParams) error {
	return s.q.UpdateUsername(ctx, UpdateUsernameParamsFromStorage(arg))
}

func (s *UserStore) UpdateEmail(ctx context.Context, arg storage.UpdateEmailParams) error {
	return s.q.UpdateEmail(ctx, UpdateEmailParamsFromStorage(arg))
}

func (s *UserStore) UpdatePassword(ctx context.Context, arg storage.UpdatePasswordParams) error {
	return s.q.UpdatePassword(ctx, UpdatePasswordParamsFromStorage(arg))
}

func (s *UserStore) UpdatePhoneNumber(ctx context.Context, arg storage.UpdatePhoneNumberParams) error {
	return s.q.UpdatePhoneNumber(ctx, UpdatePhoneNumberParamsFromStorage(arg))
}

func (s *UserStore) UpdateRole(ctx context.Context, arg storage.UpdateRoleParams) error {
	return s.q.UpdateRole(ctx, UpdateRoleParamsFromStorage(arg))
}

func (s *UserStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteUser(ctx, id.String())
}

func (s *UserStore) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeleteUser(ctx, id.String())
}

func (s *UserStore) RestoreUser(ctx context.Context, id uuid.UUID) error {
	return s.q.RestoreUser(ctx, id.String())
}

type PermissionStore struct {
	q *sqli.Queries
}

func NewPermissionStore(q *sqli.Queries) *PermissionStore {
	return &PermissionStore{q: q}
}

func (s *PermissionStore) CreatePermission(ctx context.Context, arg storage.CreatePermissionParams) error {
	return s.q.CreatePermission(ctx, CreatePermissionParamsFromStorage(arg))
}

func (s *PermissionStore) GetPermissionByID(ctx context.Context, id uuid.UUID) (*storage.DivaPermission, error) {
	p, err := s.q.GetPermissionByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return DivaPermissionToStorage(p), nil
}

func (s *PermissionStore) GetPermissionByName(ctx context.Context, name string) (*storage.DivaPermission, error) {
	p, err := s.q.GetPermissionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return DivaPermissionToStorage(p), nil
}

func (s *PermissionStore) ListPermissions(ctx context.Context, arg storage.ListPermissionsParams) ([]storage.DivaPermission, error) {
	rows, err := s.q.ListPermissions(ctx, ListPermissionsParamsFromStorage(arg))
	if err != nil {
		return nil, err
	}
	perms := make([]storage.DivaPermission, len(rows))
	for i := range rows {
		perms[i] = *DivaPermissionToStorage(rows[i])
	}
	return perms, nil
}

func (s *PermissionStore) CountPermissions(ctx context.Context) (int64, error) {
	return s.q.CountPermissions(ctx)
}

func (s *PermissionStore) UpdatePermission(ctx context.Context, arg storage.UpdatePermissionParams) error {
	return s.q.UpdatePermission(ctx, UpdatePermissionParamsFromStorage(arg))
}

func (s *PermissionStore) UpdatePermissionAction(ctx context.Context, arg storage.UpdatePermissionActionParams) error {
	return s.q.UpdatePermissionAction(ctx, UpdatePermissionActionParamsFromStorage(arg))
}

func (s *PermissionStore) UpdatePermissionRoleLevel(ctx context.Context, arg storage.UpdatePermissionRoleLevelParams) error {
	return s.q.UpdatePermissionRoleLevel(ctx, UpdatePermissionRoleLevelParamsFromStorage(arg))
}

func (s *PermissionStore) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return s.q.DeletePermission(ctx, id.String())
}

func (s *PermissionStore) SoftDeletePermission(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeletePermission(ctx, id.String())
}

func (s *PermissionStore) RestorePermission(ctx context.Context, id uuid.UUID) error {
	return s.q.RestorePermission(ctx, id.String())
}

type SessionStore struct {
	q *sqli.Queries
}

func NewSessionStore(q *sqli.Queries) *SessionStore {
	return &SessionStore{q: q}
}

func (s *SessionStore) CreateSession(ctx context.Context, arg storage.CreateSessionParams) error {
	return s.q.CreateSession(ctx, CreateSessionParamsFromStorage(arg))
}

func (s *SessionStore) GetSessionByID(ctx context.Context, id uuid.UUID) (*storage.DivaSession, error) {
	ss, err := s.q.GetSessionByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return DivaSessionToStorage(ss), nil
}

func (s *SessionStore) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]storage.DivaSession, error) {
	rows, err := s.q.ListSessionsByUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	sessions := make([]storage.DivaSession, len(rows))
	for i := range rows {
		sessions[i] = *DivaSessionToStorage(rows[i])
	}
	return sessions, nil
}

func (s *SessionStore) UpdateSession(ctx context.Context, arg storage.UpdateSessionParams) error {
	return s.q.UpdateSession(ctx, UpdateSessionParamsFromStorage(arg))
}

func (s *SessionStore) UpdateSessionStatus(ctx context.Context, arg storage.UpdateSessionStatusParams) error {
	return s.q.UpdateSessionStatus(ctx, UpdateSessionStatusParamsFromStorage(arg))
}

func (s *SessionStore) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteSession(ctx, id.String())
}

func (s *SessionStore) DeleteSessionsByUser(ctx context.Context, userID uuid.UUID) error {
	return s.q.DeleteSessionsByUser(ctx, userID.String())
}

func (s *SessionStore) DeleteExpiredSessions(ctx context.Context) error {
	return s.q.DeleteExpiredSessions(ctx)
}

type UserStateStore struct {
	q *sqli.Queries
}

func NewUserStateStore(q *sqli.Queries) *UserStateStore {
	return &UserStateStore{q: q}
}

func (s *UserStateStore) CreateUserState(ctx context.Context, arg storage.CreateUserStateParams) error {
	return s.q.CreateUserState(ctx, CreateUserStateParamsFromStorage(arg))
}

func (s *UserStateStore) GetUserStateByUserID(ctx context.Context, userID uuid.UUID) (*storage.DivaUserState, error) {
	us, err := s.q.GetUserStateByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	return DivaUserStateToStorage(us), nil
}

func (s *UserStateStore) UpdateLastActiveAt(ctx context.Context, userID uuid.UUID) error {
	return s.q.UpdateLastActiveAt(ctx, userID.String())
}

func (s *UserStateStore) UpdateUserStatus(ctx context.Context, arg storage.UpdateUserStatusParams) error {
	return s.q.UpdateUserStatus(ctx, UpdateUserStatusParamsFromStorage(arg))
}

func (s *UserStateStore) UpdateUserVerified(ctx context.Context, arg storage.UpdateUserVerifiedParams) error {
	return s.q.UpdateUserVerified(ctx, UpdateUserVerifiedParamsFromStorage(arg))
}

type UserProfileStore struct {
	q *sqli.Queries
}

func NewUserProfileStore(q *sqli.Queries) *UserProfileStore {
	return &UserProfileStore{q: q}
}

func (s *UserProfileStore) CreateUserProfile(ctx context.Context, arg storage.CreateUserProfileParams) error {
	return s.q.CreateUserProfile(ctx, CreateUserProfileParamsFromStorage(arg))
}

func (s *UserProfileStore) GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*storage.DivaUserProfile, error) {
	p, err := s.q.GetUserProfileByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	return DivaUserProfileToStorage(p), nil
}

func (s *UserProfileStore) UpdateUserProfile(ctx context.Context, arg storage.UpdateUserProfileParams) error {
	return s.q.UpdateUserProfile(ctx, UpdateUserProfileParamsFromStorage(arg))
}

func (s *UserProfileStore) UpdateUserProfileAvatar(ctx context.Context, arg storage.UpdateUserProfileAvatarParams) error {
	return s.q.UpdateUserProfileAvatar(ctx, UpdateUserProfileAvatarParamsFromStorage(arg))
}

type UserPreferenceStore struct {
	q *sqli.Queries
}

func NewUserPreferenceStore(q *sqli.Queries) *UserPreferenceStore {
	return &UserPreferenceStore{q: q}
}

func (s *UserPreferenceStore) CreateUserPreferences(ctx context.Context, arg storage.CreateUserPreferencesParams) error {
	return s.q.CreateUserPreferences(ctx, CreateUserPreferencesParamsFromStorage(arg))
}

func (s *UserPreferenceStore) GetPreferencesByID(ctx context.Context, id uuid.UUID) (*storage.DivaUserPreference, error) {
	p, err := s.q.GetPreferencesByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return DivaUserPreferenceToStorage(p), nil
}

func (s *UserPreferenceStore) GetPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]storage.DivaUserPreference, error) {
	rows, err := s.q.GetPreferencesByUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	prefs := make([]storage.DivaUserPreference, len(rows))
	for i := range rows {
		prefs[i] = *DivaUserPreferenceToStorage(rows[i])
	}
	return prefs, nil
}

func (s *UserPreferenceStore) UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) error {
	return s.q.UpdateUserPreferences(ctx, UpdateUserPreferencesParamsFromStorage(arg))
}

type UserPermissionStore struct {
	q *sqli.Queries
}

func NewUserPermissionStore(q *sqli.Queries) *UserPermissionStore {
	return &UserPermissionStore{q: q}
}

func (s *UserPermissionStore) CreateUserPermission(ctx context.Context, arg storage.CreateUserPermissionParams) error {
	return s.q.CreateUserPermission(ctx, CreateUserPermissionParamsFromStorage(arg))
}

func (s *UserPermissionStore) GetUserPermission(ctx context.Context, arg storage.GetUserPermissionParams) (*storage.DivaUserPermission, error) {
	p, err := s.q.GetUserPermission(ctx, GetUserPermissionParamsFromStorage(arg))
	if err != nil {
		return nil, err
	}
	return DivaUserPermissionToStorage(p), nil
}

func (s *UserPermissionStore) GetUserPermissionByName(ctx context.Context, arg storage.GetUserPermissionByNameParams) (*storage.DivaUserPermission, error) {
	p, err := s.q.GetUserPermissionByName(ctx, GetUserPermissionByNameParamsFromStorage(arg))
	if err != nil {
		return nil, err
	}
	return DivaUserPermissionToStorage(p), nil
}

func (s *UserPermissionStore) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]storage.DivaUserPermission, error) {
	rows, err := s.q.GetUserPermissions(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	perms := make([]storage.DivaUserPermission, len(rows))
	for i := range rows {
		perms[i] = *DivaUserPermissionToStorage(rows[i])
	}
	return perms, nil
}

func (s *UserPermissionStore) UpdateUserPermission(ctx context.Context, arg storage.UpdateUserPermissionParams) error {
	return s.q.UpdateUserPermission(ctx, UpdateUserPermissionParamsFromStorage(arg))
}

func (s *UserPermissionStore) DeleteUserPermission(ctx context.Context, arg storage.DeleteUserPermissionParams) error {
	return s.q.DeleteUserPermission(ctx, DeleteUserPermissionParamsFromStorage(arg))
}

type UserActionStore struct {
	q *sqli.Queries
}

func NewUserActionStore(q *sqli.Queries) *UserActionStore {
	return &UserActionStore{q: q}
}

func (s *UserActionStore) CreateUserAction(ctx context.Context, arg storage.CreateUserActionParams) error {
	return s.q.CreateUserAction(ctx, CreateUserActionParamsFromStorage(arg))
}

func (s *UserActionStore) GetUserActionByID(ctx context.Context, id uuid.UUID) (*storage.DivaAction, error) {
	a, err := s.q.GetUserActionByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return DivaActionToStorage(a), nil
}

func (s *UserActionStore) GetUserActionByUserAndName(ctx context.Context, arg storage.GetUserActionByUserAndNameParams) (*storage.DivaAction, error) {
	a, err := s.q.GetUserActionByUserAndName(ctx, GetUserActionByUserAndNameParamsFromStorage(arg))
	if err != nil {
		return nil, err
	}
	return DivaActionToStorage(a), nil
}

func (s *UserActionStore) ListActionsByUser(ctx context.Context, userID uuid.UUID) ([]storage.DivaAction, error) {
	rows, err := s.q.ListActionsByUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	actions := make([]storage.DivaAction, len(rows))
	for i := range rows {
		actions[i] = *DivaActionToStorage(rows[i])
	}
	return actions, nil
}

func (s *UserActionStore) DeleteUserAction(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteUserAction(ctx, id.String())
}

func (s *UserActionStore) DeleteUserActionByUser(ctx context.Context, userID uuid.UUID) error {
	return s.q.DeleteUserActionByUser(ctx, userID.String())
}

type UserVerificationStore struct {
	q *sqli.Queries
}

func NewUserVerificationStore(q *sqli.Queries) *UserVerificationStore {
	return &UserVerificationStore{q: q}
}

func (s *UserVerificationStore) CreateUserVerification(ctx context.Context, arg storage.CreateUserVerificationParams) error {
	return s.q.CreateUserVerification(ctx, CreateUserVerificationParamsFromStorage(arg))
}

func (s *UserVerificationStore) GetUserVerification(ctx context.Context, actionID uuid.UUID) (*storage.DivaActionVerification, error) {
	v, err := s.q.GetUserVerification(ctx, actionID.String())
	if err != nil {
		return nil, err
	}
	return DivaActionVerificationToStorage(v), nil
}

func (s *UserVerificationStore) UpdateUserVerification(ctx context.Context, arg storage.UpdateUserVerificationParams) error {
	return s.q.UpdateUserVerification(ctx, UpdateUserVerificationParamsFromStorage(arg))
}

func (s *UserVerificationStore) DeleteUserVerification(ctx context.Context, actionID uuid.UUID) error {
	return s.q.DeleteUserVerification(ctx, actionID.String())
}
