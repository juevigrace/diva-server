package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/api/core/user/permissions"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserProfileRepo struct {
	queries   *db.Queries
	upRepo *permissions.UserPermissionRepo
}

func NewUserProfileRepo(
	queries *db.Queries,
	upRepo *permissions.UserPermissionRepo,
) *UserProfileRepo {
	return &UserProfileRepo{
		queries:   queries,
		upRepo: upRepo,
	}
}

func (s *UserProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	row, err := s.queries.GetUserProfileByUserID(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		return nil, err
	}
	return models.UserProfileFromDB(&row), nil
}

func (s *UserProfileRepo) Create(ctx context.Context, session *models.Session, uid uuid.UUID, dto *dtos.CreateProfileDto) error {
	profile := &models.UserProfile{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Bio:       dto.Bio,
	}
	if err := s.queries.CreateUserProfile(ctx, *profile.DBCreate(uid)); err != nil {
		return err
	}
	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_PROFILE_WRITE]; ok {
			if err := s.upRepo.Delete(ctx, session.User.ID, perm.Permission.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *UserProfileRepo) Update(ctx context.Context, userID uuid.UUID, dto *dtos.UpdateProfileDto) error {
	profile := &models.UserProfile{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate,
		Alias:     dto.Alias,
		Bio:       dto.Bio,
	}
	return s.queries.UpdateUserProfile(ctx, *profile.DBUpdate(userID))
}

func (s *UserProfileRepo) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	return s.queries.UpdateUserProfileAvatar(ctx, db.UpdateUserProfileAvatarParams{
		Avatar: avatar,
		UserID: models.UUIDPtrToDB(&userID),
	})
}
