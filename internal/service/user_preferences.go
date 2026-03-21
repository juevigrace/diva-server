package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/storage/db"
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

	params := &db.CreateUserPreferencesParams{
		ID:                  pgtype.UUID{Bytes: id, Valid: true},
		UserID:              pgtype.UUID{Bytes: userID, Valid: true},
		Theme:               models.ThemeFromString(dto.Theme).ToDB(),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
		CreatedAt:           models.ToTimestamptzPtr(&dto.CreatedAt),
		UpdatedAt:           models.ToTimestamptzPtr(&dto.UpdatedAt),
	}

	return s.repo.Create(ctx, params)
}

func (s *UserPreferencesService) Update(ctx context.Context, dto *dtos.UserPreferencesDto) error {
	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return err
	}

	params := &db.UpdateUserPreferencesParams{
		Theme:               models.ThemeFromString(dto.Theme).ToDB(),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
		UpdatedAt:           models.ToTimestamptzPtr(&dto.UpdatedAt),
		ID:                  pgtype.UUID{Bytes: id, Valid: true},
	}
	return s.repo.Update(ctx, params)
}
