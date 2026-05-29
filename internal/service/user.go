package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
)

type UserService struct {
	repo       *repo.UserRepo
	uaService  *UserActionsService
	upService  *UserPermissionService
	uprService *UserProfileService
	uvService  *UserVerificationService
}

func NewUserService(
	repo *repo.UserRepo,
	uaService *UserActionsService,
	upService *UserPermissionService,
	uprService *UserProfileService,
	uvService *UserVerificationService,
) *UserService {
	return &UserService{
		repo:       repo,
		uaService:  uaService,
		uvService:  uvService,
		upService:  upService,
		uprService: uprService,
	}
}

func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *UserService) GetAll(ctx context.Context, pagination *models.Pagination) ([]models.User, error) {
	return s.repo.ListUsers(ctx, pagination.GetLimit(), pagination.GetOffset())
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	dbUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

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
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) CheckUsernameAvailable(ctx context.Context, username string) (bool, error) {
	_, err := s.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(models.ErrUserNotFound, err) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) CheckEmailAvailable(ctx context.Context, email string) (bool, error) {
	_, err := s.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(models.ErrUserNotFound, err) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (s *UserService) GetByUsernameOrEmail(ctx context.Context, value string) (*models.User, error) {
	return s.repo.GetByUsernameOrEmail(ctx, value)
}

func (s *UserService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	_, err := s.repo.GetByUsername(ctx, username)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, models.ErrUserNotFound) {
		return true, nil
	}
	return false, err
}

func (s *UserService) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, models.ErrUserNotFound) {
		return true, nil
	}
	return false, err
}

func (s *UserService) Create(ctx context.Context, dto *dtos.CreateUserDto) (uuid.UUID, error) {
	id := uuid.New()

	passwordHash, err := util.HashPassword(dto.Password)
	if err != nil {
		return uuid.Nil, err
	}

	user := &models.User{
		ID:           id,
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Verified:     false,
		Role:         models.ROLE_USER,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return uuid.Nil, err
	}

	// TODO: need to create any other user related data here

	if _, err := s.uaService.Create(ctx, id, models.ActionUserVerification); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *UserService) VerifyUser(ctx context.Context, actionID uuid.UUID) error {
	dbAction, err := s.uvService.GetOneById(ctx, actionID)
	if err != nil {
		return err
	}

	if !dbAction.Verified {
		return models.ErrActionNotVerified
	}

	if err := s.repo.UpdateVerified(ctx, true, dbAction.Action.UserID); err != nil {
		return err
	}

	if err := s.uaService.Delete(ctx, dbAction.Action.ID); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdatePassword(ctx context.Context, uid uuid.UUID, newPassword string) error {
	dbUser, err := s.GetByID(ctx, uid)
	if err != nil {
		return err
	}

	if util.ValidatePassword(newPassword, dbUser.PasswordHash) {
		return models.ErrSamePassword
	}

	newHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, newHash, uid); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdatePhoneNumber(ctx context.Context, phone string, userID uuid.UUID) error {
	return s.repo.UpdatePhoneNumber(ctx, phone, userID)
}

func (s *UserService) UpdateUsername(ctx context.Context, username string, userID uuid.UUID) error {
	return s.repo.UpdateUsername(ctx, username, userID)
}

func (s *UserService) UpdateEmail(ctx context.Context, email string, userID uuid.UUID) error {
	return s.repo.UpdateEmail(ctx, email, userID)
}

func (s *UserService) UpdateVerified(ctx context.Context, verified bool, userID uuid.UUID) error {
	return s.repo.UpdateVerified(ctx, verified, userID)
}

func (s *UserService) UpdateRole(ctx context.Context, role models.Role, userID uuid.UUID) error {
	return s.repo.UpdateRole(ctx, role, userID)
}

func (s *UserService) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	return s.repo.SoftDelete(ctx, userID)
}

func (s *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Delete(ctx, userID)
}
