package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPreferencesRepository struct {
	queries *db.Queries
}

func NewUserPreferencesRepository(queries *db.Queries) *UserPreferencesRepository {
	return &UserPreferencesRepository{queries: queries}
}

func (r *UserPreferencesRepository) Create(ctx context.Context, params *models.UserPreferences) error {
	return r.queries.CreateUserPreferences(ctx, db.CreateUserPreferencesParams{
		ID:                  pgtype.UUID{Bytes: params.ID, Valid: true},
		UserID:              pgtype.UUID{Bytes: params.UserID, Valid: true},
		Theme:               params.Theme.ToDB(),
		OnboardingCompleted: params.OnboardingCompleted,
		Language:            params.Language,
		CreatedAt:           pgtype.Timestamptz{Time: time.UnixMilli(params.CreatedAt), Valid: true},
		UpdatedAt:           pgtype.Timestamptz{Time: time.UnixMilli(params.UpdatedAt), Valid: true},
	})
}

func (r *UserPreferencesRepository) Update(ctx context.Context, params *models.UserPreferences) error {
	return r.queries.UpdateUserPreferences(ctx, db.UpdateUserPreferencesParams{
		Theme:     params.Theme.ToDB(),
		Language:  params.Language,
		UpdatedAt: pgtype.Timestamptz{Time: time.UnixMilli(params.UpdatedAt), Valid: true},
		ID:        pgtype.UUID{Bytes: params.ID, Valid: true},
	})
}

func (r *UserPreferencesRepository) CreateBatch(ctx context.Context, params []*models.UserPreferences) error {
	for _, p := range params {
		if err := r.Create(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserPreferencesRepository) UpdateBatch(ctx context.Context, params []*models.UserPreferences) error {
	for _, p := range params {
		if err := r.Update(ctx, p); err != nil {
			return err
		}
	}
	return nil
}
