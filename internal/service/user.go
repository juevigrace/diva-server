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
	repo      *repo.UserRepository
	uaService *UserActionsService
}

func NewUserService(repo *repo.UserRepository, uaService *UserActionsService) *UserService {
	return &UserService{repo: repo, uaService: uaService}
}

func (s *UserService) Create(ctx context.Context, dto *dtos.CreateUserDto) (*uuid.UUID, error) {
	id := uuid.New()

	passwordHash, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           id,
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Alias:        dto.Alias,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// TODO: need to create any other user related data here

	if _, err := s.uaService.Create(ctx, models.ActionUserVerification, &id); err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, dto *dtos.UpdateProfileDto) error {
	user := &models.User{
		ID:        userID,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Avatar:    dto.Avatar,
		Bio:       dto.Bio,
	}

	return s.repo.UpdateProfile(ctx, user)
}

func (s *UserService) UpdatePassword(ctx context.Context, session *models.Session, newPassword string) error {
	if util.ValidatePassword(newPassword, session.User.PasswordHash) {
		return models.ErrPasswordsMatch
	}

	if util.ValidatePassword(newPassword, session.User.PasswordHash) {
		return models.ErrSamePassword
	}

	newHash, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, newHash, &session.User.ID)
}

func (s *UserService) VerifyUser(ctx context.Context, userID *uuid.UUID) error {
	if err := s.repo.VerifyUser(ctx, userID); err != nil {
		return err
	}

	if err := s.uaService.Delete(ctx, &models.UserAction{
		UserID: *userID,
		Action: models.ActionUserVerification,
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserService) CreateBatch(ctx context.Context, params []*models.User) error {
	return s.repo.CreateBatch(ctx, params)
}

func (s *UserService) UpdatePhoneNumber(ctx context.Context, phone string, userID *uuid.UUID) error {
	return s.repo.UpdatePhoneNumber(ctx, phone, userID)
}

func (s *UserService) UpdateUsername(ctx context.Context, username string, userID *uuid.UUID) error {
	return s.repo.UpdateUsername(ctx, username, userID)
}

func (s *UserService) UpdateEmail(ctx context.Context, dto *dtos.UpdateEmailDto, userID *uuid.UUID) error {
	return s.repo.UpdateEmail(ctx, dto, userID)
}

func (s *UserService) Delete(ctx context.Context, userID *uuid.UUID) error {
	return s.repo.Delete(ctx, userID)
}

func (s *UserService) CheckUsernameAvailable(ctx context.Context, username string) (bool, error) {
	_, err := s.repo.GetByUsername(ctx, username)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, models.ErrUserNotFound) {
		return true, nil
	}
	return false, err
}

func (s *UserService) CheckEmailAvailable(ctx context.Context, email string) (bool, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, models.ErrUserNotFound) {
		return true, nil
	}
	return false, err
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) GetByUsernameOrEmail(ctx context.Context, value string) (*models.User, error) {
	return s.repo.GetByUsernameOrEmail(ctx, value)
}

func (s *UserService) GetByID(ctx context.Context, userID *uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *UserService) GetAll(ctx context.Context, pagination *models.Pagination) ([]*models.User, error) {
	return s.repo.GetAll(ctx, pagination.GetLimit(), pagination.GetOffset())
}

func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
