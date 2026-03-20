package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPreferencesRepository struct {
	queries *db.Queries
}

func NewUserPreferencesRepository(queries *db.Queries) *UserPreferencesRepository {
	return &UserPreferencesRepository{queries: queries}
}

func (r *UserPreferencesRepository) Create(ctx context.Context, userID uuid.UUID, dto *dtos.UserPreferencesDto) error {
	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().UnixMilli()
	onboardingCompleted := false
	if dto.OnboardingCompleted != nil {
		onboardingCompleted = *dto.OnboardingCompleted
	}
	params := db.CreateUserPreferencesParams{
		ID:                  pgtype.UUID{Bytes: id, Valid: true},
		UserID:              pgtype.UUID{Bytes: userID, Valid: true},
		Theme:               models.ThemeFromString(dto.Theme).ToDB(),
		OnboardingCompleted: onboardingCompleted,
		Language:            dto.Language,
		CreatedAt:           models.ToTimestamptzPtr(&now),
		UpdatedAt:           models.ToTimestamptzPtr(&now),
	}
	return r.queries.CreateUserPreferences(ctx, params)
}

func (r *UserPreferencesRepository) Update(ctx context.Context, dto *dtos.UserPreferencesDto) error {
	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return err
	}
	now := time.Now().UTC().UnixMilli()
	onboardingCompleted := false
	if dto.OnboardingCompleted != nil {
		onboardingCompleted = *dto.OnboardingCompleted
	}
	params := db.UpdateUserPreferencesParams{
		Theme:               models.ThemeFromString(dto.Theme).ToDB(),
		OnboardingCompleted: onboardingCompleted,
		Language:            dto.Language,
		UpdatedAt:           models.ToTimestamptzPtr(&now),
		ID:                  pgtype.UUID{Bytes: id, Valid: true},
	}
	return r.queries.UpdateUserPreferences(ctx, params)
}
