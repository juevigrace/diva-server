package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/api/core/session"
	"github.com/juevigrace/diva-server/internal/api/core/user/actions"
	"github.com/juevigrace/diva-server/internal/api/core/user/permissions"
	"github.com/juevigrace/diva-server/internal/api/core/user/profile"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/pkg/bcrypt"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/storage"
)

type UserRepo struct {
	store   storage.UserStore
	sRepo   *session.SessionRepo
	uaRepo  *actions.UserActionsRepo
	upRepo  *permissions.UserPermissionRepo
	uproRepo *profile.UserProfileRepo
	usRepo  *UserStateRepo
}

func NewUserRepo(
	store storage.UserStore,
	sRepo *session.SessionRepo,
	uaRepo *actions.UserActionsRepo,
	upRepo *permissions.UserPermissionRepo,
	uproRepo *profile.UserProfileRepo,
	usRepo *UserStateRepo,
) *UserRepo {
	return &UserRepo{
		store:   store,
		sRepo:   sRepo,
		uaRepo:  uaRepo,
		upRepo:  upRepo,
		uproRepo: uproRepo,
		usRepo:  usRepo,
	}
}

func (s *UserRepo) Count(ctx context.Context) (int64, error) {
	return s.store.CountUsers(ctx)
}

func (s *UserRepo) GetAll(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	rows, err := s.store.ListUsers(ctx, &storage.ListUsersParams{
		Limit:  int64(pagination.GetLimit()),
		Offset: int64(pagination.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(rows))
	for i := range rows {
		users[i] = *models.UserFromDB(&rows[i])
	}

	return users, nil
}

func (s *UserRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	row, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	dbUser := models.UserFromDB(row)

	dbProfile, err := s.uproRepo.GetByUserID(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}
	dbUser.Profile = dbProfile

	dbPerms, err := s.upRepo.GetByUser(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}

	dbUser.Permissions = make(map[models.PermissionAction]models.UserPermission, len(dbPerms))
	for _, perm := range dbPerms {
		dbUser.Permissions[perm.Permission.Action] = perm
	}

	dbState, err := s.usRepo.GetByUserID(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}
	dbUser.State = dbState

	return dbUser, nil
}

func (s *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return models.UserFromDB(row), nil
}

func (s *UserRepo) CheckUsernameAvailable(ctx context.Context, username string) (bool, error) {
	_, err := s.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (s *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return models.UserFromDB(row), nil
}

func (s *UserRepo) CheckEmailAvailable(ctx context.Context, email string) (bool, error) {
	_, err := s.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (s *UserRepo) GetByUsernameOrEmail(ctx context.Context, value string) (*models.User, error) {
	row, err := s.store.GetUserByUsernameOrEmail(ctx, value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return models.UserFromDB(row), nil
}

func (s *UserRepo) Create(ctx context.Context, dto *dtos.CreateUserDto) (uuid.UUID, error) {
	passwordHash, err := bcrypt.HashPassword(dto.Password)
	if err != nil {
		return uuid.Nil, err
	}

	uid := uuid.New()
	params := &models.User{
		ID:           uid,
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Role:         models.ROLE_USER,
	}

	if err := s.store.CreateUser(ctx, params.DBCreate()); err != nil {
		return uuid.Nil, err
	}

	if err := s.usRepo.Create(ctx, params.ID, &models.UserState{}); err != nil {
		if err := s.Delete(ctx, params.ID); err != nil {
			return uuid.Nil, err
		}
		return uuid.Nil, err
	}

	if _, err := s.uaRepo.Create(ctx, params.ID, models.ActionUserVerification); err != nil {
		return uuid.Nil, err
	}

	perms := []models.PermissionAction{models.PERMISSION_USERS_PROFILE_WRITE, models.PERMISSION_USERS_PREFERENCES_WRITE}
	for i := range perms {
		if err := s.upRepo.CreateByName(ctx, perms[i], nil, true, nil, uid); err != nil {
			if err := s.Delete(ctx, uid); err != nil {
				return uuid.Nil, err
			}
			return uuid.Nil, err
		}
	}

	return uid, nil
}

func (s *UserRepo) UpdatePasswordConfirm(ctx context.Context, aid, uid uuid.UUID) error {
	exp := time.Now().UTC().Add(15 * time.Minute).UnixMilli()
	if err := s.upRepo.CreateByName(ctx, models.PERMISSION_USERS_PASSWORD_WRITE, nil, true, &exp, uid); err != nil {
		return err
	}

	if err := s.uaRepo.Delete(ctx, aid); err != nil {
		return err
	}

	return nil
}

func (s *UserRepo) UpdatePassword(ctx context.Context, session *models.Session, uid uuid.UUID, newPassword string) error {
	dbUser, err := s.GetByID(ctx, uid)
	if err != nil {
		return err
	}

	if bcrypt.ValidatePassword(newPassword, dbUser.PasswordHash) {
		return errs.ErrSamePassword
	}

	newHash, err := bcrypt.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.store.UpdatePassword(ctx, &storage.UpdatePasswordParams{
		PasswordHash: newHash,
		ID:           uid,
	}); err != nil {
		return err
	}

	if err := s.sRepo.Close(ctx, uid); err != nil {
		return err
	}

	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_PASSWORD_WRITE]; ok {
			if err := s.upRepo.Delete(ctx, session.User.ID, perm.Permission.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *UserRepo) UpdatePhoneNumber(ctx context.Context, session *models.Session, phone string, uid uuid.UUID) error {
	if err := s.store.UpdatePhoneNumber(ctx, &storage.UpdatePhoneNumberParams{
		PhoneNumber: phone,
		ID:          uid,
	}); err != nil {
		return err
	}

	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_PHONE_WRITE]; ok {
			if err := s.upRepo.Delete(ctx, session.User.ID, perm.Permission.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *UserRepo) UpdateUsername(ctx context.Context, session *models.Session, username string, uid uuid.UUID) error {
	if err := s.store.UpdateUsername(ctx, &storage.UpdateUsernameParams{
		Username: username,
		ID:       uid,
	}); err != nil {
		return err
	}

	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_USERNAME_WRITE]; ok {
			if err := s.upRepo.Delete(ctx, session.User.ID, perm.Permission.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *UserRepo) UpdateEmail(ctx context.Context, session *models.Session, email string, uid uuid.UUID) error {
	if err := s.store.UpdateEmail(ctx, &storage.UpdateEmailParams{
		Email: email,
		ID:    uid,
	}); err != nil {
		return err
	}

	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_EMAIL_WRITE]; ok {
			if err := s.upRepo.Delete(ctx, session.User.ID, perm.Permission.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *UserRepo) UpdateRole(ctx context.Context, role models.Role, userID uuid.UUID) error {
	return s.store.UpdateRole(ctx, &storage.UpdateRoleParams{
		Role: role.ToDB(),
		ID:   userID,
	})
}

func (s *UserRepo) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	if err := s.store.SoftDeleteUser(ctx, userID); err != nil {
		return err
	}

	return s.usRepo.UpdateStatus(ctx, models.USER_STATUS_INACTIVE, userID)
}

func (s *UserRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.store.DeleteUser(ctx, userID)
}

func (r *UserRepo) Restore(ctx context.Context, userID uuid.UUID) error {
	if err := r.store.RestoreUser(ctx, userID); err != nil {
		return err
	}

	return r.usRepo.UpdateStatus(ctx, models.USER_STATUS_ACTIVE, userID)
}
