package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserProfileService struct {
	repo *repo.UserProfileRepo
}

func NewUserProfileService(repo *repo.UserProfileRepo) *UserProfileService {
	return &UserProfileService{repo: repo}
}

func (s *UserProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *UserProfileService) Create(ctx context.Context, userID uuid.UUID, dto *dtos.CreateProfileDto) error {
	profile := &models.UserProfile{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Bio:       dto.Bio,
	}

	return s.repo.Create(ctx, userID, profile)
}

func (s *UserProfileService) Update(ctx context.Context, userID uuid.UUID, dto *dtos.UpdateProfileDto) error {
	profile := &models.UserProfile{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Bio:       dto.Bio,
	}

	return s.repo.Update(ctx, userID, profile)
}

func (s *UserProfileService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	// TODO: handle here file of the image
	return s.repo.UpdateAvatar(ctx, userID, avatar)
}
