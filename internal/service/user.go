package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserService struct {
	repo         *repo.UserRepository
	verification *VerificationService
}

func NewUserService(repo *repo.UserRepository, verification *VerificationService) *UserService {
	return &UserService{repo: repo, verification: verification}
}

func (s *UserService) Create(ctx context.Context, dto *dtos.CreateUserDto) (uuid.UUID, error) {
	id := uuid.New()

	passwordHash, err := util.HashPassword(dto.Password)
	if err != nil {
		return uuid.Nil, err
	}

	params := &db.CreateUserParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Alias:        dto.Alias,
	}

	return s.repo.Create(ctx, params)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, dto *dtos.UpdateProfileDto) error {
	params := db.UpdateProfileParams{
		Alias:  dto.Alias,
		Avatar: dto.Avatar,
		Bio:    dto.Bio,
		ID:     pgtype.UUID{Bytes: userID, Valid: true},
	}
	return s.repo.UpdateProfile(ctx, &params)
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

	params := &db.UpdatePasswordParams{
		PasswordHash: newHash,
		ID:           pgtype.UUID{Bytes: session.User.ID, Valid: true},
	}

	return s.repo.UpdatePassword(ctx, params)
}

func (s *UserService) UpdateVerified(ctx context.Context, userID *uuid.UUID) error {
	return s.repo.VerifyUser(ctx, userID)
}

func (s *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
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

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *UserService) GetAll(ctx context.Context, pagination *models.Pagination) ([]*models.User, error) {
	return s.repo.GetAll(ctx, pagination.GetLimit(), pagination.GetOffset())
}

func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
