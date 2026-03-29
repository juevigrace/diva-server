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
	repo *repo.UserPreferencesRepository
}

func NewUserPreferencesService(repo *repo.UserPreferencesRepository) *UserPreferencesService {
	return &UserPreferencesService{repo: repo}
}

func (s *UserPreferencesService) Create(ctx context.Context, userID uuid.UUID, dto *dtos.UserPreferencesDto) error {
	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return err
	}

	pref := &models.UserPreferences{
		ID:                  id,
		UserID:              userID,
		Theme:               models.ThemeFromString(dto.Theme),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
		LastSyncAt:          time.Now().UTC().UnixMilli(),
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
	}

	return s.repo.Create(ctx, pref)
}

func (s *UserPreferencesService) Update(ctx context.Context, dto *dtos.UserPreferencesDto) error {
	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return err
	}

	pref := &models.UserPreferences{
		ID:         id,
		Theme:      models.ThemeFromString(dto.Theme),
		Language:   dto.Language,
		LastSyncAt: time.Now().UTC().UnixMilli(),
		UpdatedAt:  dto.UpdatedAt,
	}

	return s.repo.Update(ctx, pref)
}

func (s *UserPreferencesService) CreateBatch(ctx context.Context, params []*models.UserPreferences) error {
	return s.repo.CreateBatch(ctx, params)
}

func (s *UserPreferencesService) UpdateBatch(ctx context.Context, params []*models.UserPreferences) error {
	return s.repo.UpdateBatch(ctx, params)
}
