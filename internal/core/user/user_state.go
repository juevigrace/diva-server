package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserStateService struct {
	queries *db.Queries
}

func NewUserStateService(queries *db.Queries) *UserStateService {
	return &UserStateService{
		queries: queries,
	}
}

func (s *UserStateService) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserState, error) {
	row, err := s.queries.GetUserStateByUserID(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		return nil, err
	}
	return models.UserStateFromDB(&row), nil
}

func (s *UserStateService) Create(ctx context.Context, userID uuid.UUID, us *models.UserState) error {
	return s.queries.CreateUserState(ctx, *us.DBCreate(userID))
}

func (s *UserStateService) UpdateVerified(ctx context.Context, verified bool, userID uuid.UUID) error {
	return s.queries.UpdateUserVerified(ctx, db.UpdateUserVerifiedParams{
		Verified: verified,
		UserID:   models.UUIDPtrToDB(&userID),
	})
}

func (s *UserStateService) UpdateStatus(ctx context.Context, status models.UserStatus, userID uuid.UUID) error {
	return s.queries.UpdateUserStatus(ctx, db.UpdateUserStatusParams{
		Status: status.ToDB(),
		UserID: models.UUIDPtrToDB(&userID),
	})
}

func (s *UserStateService) UpdateLastActiveAt(ctx context.Context, userID uuid.UUID) error {
	return s.queries.UpdateLastActiveAt(ctx, models.UUIDPtrToDB(&userID))
}
