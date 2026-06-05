package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserService struct {
	queries    *db.Queries
	uaService  *UserActionsService
	upService  *UserPermissionService
	uprService *UserProfileService
	uvService  *UserVerificationService
}

func NewUserService(
	queries *db.Queries,
	uaService *UserActionsService,
	upService *UserPermissionService,
	uprService *UserProfileService,
	uvService *UserVerificationService,
) *UserService {
	return &UserService{
		uaService:  uaService,
		uvService:  uvService,
		upService:  upService,
		uprService: uprService,
		queries:    queries,
	}
}

func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.queries.CountUsers(ctx)
}

func (s *UserService) GetAll(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	rows, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(pagination.GetLimit()),
		Offset: int32(pagination.GetOffset()),
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

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	row, err := s.queries.GetUserByID(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	dbUser := models.UserFromDB(&row)

	dbProfile, err := s.uprService.GetByUserID(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}
	dbUser.Profile = dbProfile

	dbPerms, err := s.upService.GetByUser(ctx, dbUser.ID)
	if err != nil {
		return nil, err
	}

	dbUser.Permissions = make(map[models.PermissionAction]models.UserPermission, len(dbPerms))
	for _, perm := range dbPerms {
		dbUser.Permissions[perm.Permission.Action] = perm
	}

	return dbUser, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return models.UserFromDB(&row), nil
}

func (s *UserService) CheckUsernameAvailable(ctx context.Context, username string) (bool, error) {
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

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return models.UserFromDB(&row), nil
}

func (s *UserService) CheckEmailAvailable(ctx context.Context, email string) (bool, error) {
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

func (s *UserService) GetByUsernameOrEmail(ctx context.Context, value string) (*models.User, error) {
	row, err := s.queries.GetUserByUsernameOrEmail(ctx, value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return models.UserFromDB(&row), nil
}

func (s *UserService) Create(ctx context.Context, dto *dtos.CreateUserDto) (uuid.UUID, error) {
	passwordHash, err := util.HashPassword(dto.Password)
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

	if err := s.queries.CreateUser(ctx, *params.DBCreate()); err != nil {
		return uuid.Nil, err
	}

	if _, err := s.uaService.Create(ctx, params.ID, models.ActionUserVerification); err != nil {
		return uuid.Nil, err
	}

	perms := []models.PermissionAction{models.PERMISSION_USERS_PROFILE_WRITE, models.PERMISSION_USERS_PREFERENCES_WRITE}
	for i := range perms {
		if err := s.upService.GrantByName(ctx, perms[i], nil, nil, uid); err != nil {
			if err := s.Delete(ctx, uid); err != nil {
				return uuid.Nil, err
			}
			return uuid.Nil, err
		}
	}

	return uid, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, uid uuid.UUID, newPassword string) error {
	dbUser, err := s.GetByID(ctx, uid)
	if err != nil {
		return err
	}

	if util.ValidatePassword(newPassword, dbUser.PasswordHash) {
		return errs.ErrSamePassword
	}

	newHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		PasswordHash: newHash,
		ID:           models.UUIDPtrToDB(&uid),
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdatePhoneNumber(ctx context.Context, phone string, userID uuid.UUID) error {
	return s.queries.UpdatePhoneNumber(ctx, db.UpdatePhoneNumberParams{
		PhoneNumber: phone,
		ID:          models.UUIDPtrToDB(&userID),
	})
}

func (s *UserService) UpdateUsername(ctx context.Context, username string, userID uuid.UUID) error {
	return s.queries.UpdateUsername(ctx, db.UpdateUsernameParams{
		Username: username,
		ID:       models.UUIDPtrToDB(&userID),
	})

}

func (s *UserService) UpdateEmail(ctx context.Context, email string, userID uuid.UUID) error {
	return s.queries.UpdateEmail(ctx, db.UpdateEmailParams{
		Email: email,
		ID:    models.UUIDPtrToDB(&userID),
	})
}

func (s *UserService) UpdateVerified(ctx context.Context, verified bool, userID uuid.UUID) error {
	return s.queries.UpdateVerified(ctx, db.UpdateVerifiedParams{
		Verified: verified,
		ID:       models.UUIDPtrToDB(&userID),
	})
}

func (s *UserService) UpdateRole(ctx context.Context, role models.Role, userID uuid.UUID) error {
	return s.queries.UpdateRole(ctx, db.UpdateRoleParams{
		Role: role.ToDB(),
		ID:   models.UUIDPtrToDB(&userID),
	})
}

func (s *UserService) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	return s.queries.DeleteUser(ctx, models.UUIDPtrToDB(&userID))
}

func (s *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.queries.SoftDeleteUser(ctx, models.UUIDPtrToDB(&userID))
}

func (r *UserService) Restore(ctx context.Context, userID uuid.UUID) error {
	return r.queries.RestoreUser(ctx, models.UUIDPtrToDB(&userID))
}
