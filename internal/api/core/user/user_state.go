package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage"
)

type UserStateRepo struct {
	store storage.UserStateStore
}

func NewUserStateRepo(store storage.UserStateStore) *UserStateRepo {
	return &UserStateRepo{
		store: store,
	}
}

func (s *UserStateRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserState, error) {
	row, err := s.store.GetUserStateByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return models.UserStateFromDB(row), nil
}

func (s *UserStateRepo) Create(ctx context.Context, userID uuid.UUID, us *models.UserState) error {
	return s.store.CreateUserState(ctx, us.DBCreate(userID))
}

func (s *UserStateRepo) UpdateVerified(ctx context.Context, verified bool, userID uuid.UUID) error {
	return s.store.UpdateUserVerified(ctx, &storage.UpdateUserVerifiedParams{
		Verified: verified,
		UserID:   userID,
	})
}

func (s *UserStateRepo) UpdateStatus(ctx context.Context, status models.UserStatus, userID uuid.UUID) error {
	return s.store.UpdateUserStatus(ctx, &storage.UpdateUserStatusParams{
		Status: status.ToDB(),
		UserID: userID,
	})
}

func (s *UserStateRepo) UpdateLastActiveAt(ctx context.Context, userID uuid.UUID) error {
	return s.store.UpdateLastActiveAt(ctx, userID)
}
