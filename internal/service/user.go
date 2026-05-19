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
	repo      *repo.UserRepo
	uaService *UserActionsService
	uvService *UserVerificationService
}

func NewUserService(repo *repo.UserRepo, uaService *UserActionsService) *UserService {
	return &UserService{repo: repo, uaService: uaService}
}

func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *UserService) GetAll(ctx context.Context, pagination *models.Pagination) ([]*models.User, error) {
	return s.repo.ListUsers(ctx, pagination.GetLimit(), pagination.GetOffset())
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
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

	// TODO: create verification for action before this

	if _, err := s.uaService.Create(ctx, models.ActionUserVerification, &id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, session *models.Session, newPassword string) error {
	if util.ValidatePassword(newPassword, session.User.PasswordHash) {
		return models.ErrSamePassword
	}

	newHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, newHash, session.User.ID)
}

func (s *UserService) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.UpdateVerified(ctx, true, userID); err != nil {
		return err
	}

	// TODO: delete action here
	// if err := s.uaService.Delete(ctx, &models.UserAction{
	// 	UserID: userID,
	// 	Action: models.ActionUserVerification,
	// }); err != nil {
	// 	return err
	// }

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

func (s *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.repo.SoftDelete(ctx, userID)
}
