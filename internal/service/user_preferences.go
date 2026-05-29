package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserPreferencesService struct {
	repo *repo.UserPreferencesRepo
}

func NewUserPreferencesService(repo *repo.UserPreferencesRepo) *UserPreferencesService {
	return &UserPreferencesService{repo: repo}
}

func (s *UserPreferencesService) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPreferences, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *UserPreferencesService) GetByID(ctx context.Context, id uuid.UUID) (*models.UserPreferences, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserPreferencesService) Create(ctx context.Context, userID uuid.UUID, dto *dtos.CreateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:                  uuid.New(),
		Device:              dto.Device,
		Theme:               models.ThemeFromString(dto.Theme),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
		LastSyncAt:          time.Now().UTC().UnixMilli(),
		CreatedAt:           time.Now().UTC().UnixMilli(),
		UpdatedAt:           time.Now().UTC().UnixMilli(),
	}

	return s.repo.Create(ctx, userID, pref)
}

func (s *UserPreferencesService) Update(ctx context.Context, id uuid.UUID, dto *dtos.UpdateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:         id,
		Theme:      models.ThemeFromString(dto.Theme),
		Language:   dto.Language,
		LastSyncAt: time.Now().UTC().UnixMilli(),
		UpdatedAt:  time.Now().UTC().UnixMilli(),
	}

	return s.repo.Update(ctx, pref)
}

func (s *UserPreferencesService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserPreferencesService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUser(ctx, userID)
}
