package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserProfileService struct {
	queries *db.Queries
}

func NewUserProfileService(queries *db.Queries) *UserProfileService {
	return &UserProfileService{
		queries: queries,
	}
}

func (s *UserProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	row, err := s.queries.GetUserProfileByUserID(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		return nil, err
	}
	return models.UserProfileFromDB(&row), nil
}

func (s *UserProfileService) Create(ctx context.Context, userID uuid.UUID, dto *dtos.CreateProfileDto) error {
	profile := &models.UserProfile{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Bio:       dto.Bio,
	}
	return s.queries.CreateUserProfile(ctx, *profile.DBCreate(userID))
}

func (s *UserProfileService) Update(ctx context.Context, userID uuid.UUID, dto *dtos.UpdateProfileDto) error {
	profile := &models.UserProfile{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Bio:       dto.Bio,
	}
	return s.queries.UpdateUserProfile(ctx, *profile.DBUpdate(userID))
}

func (s *UserProfileService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	// TODO: handle here file of the image
	return s.queries.UpdateUserProfileAvatar(ctx, db.UpdateUserProfileAvatarParams{
		Avatar: avatar,
		UserID: models.UUIDPtrToDB(&userID),
	})
}
