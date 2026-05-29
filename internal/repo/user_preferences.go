package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPreferencesRepo struct {
	queries *db.Queries
}

func NewUserPreferencesRepo(queries *db.Queries) *UserPreferencesRepo {
	return &UserPreferencesRepo{queries: queries}
}

func (r *UserPreferencesRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPreferences, error) {
	rows, err := r.queries.GetPreferencesByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	prefs := make([]*models.UserPreferences, len(rows))
	for i := range rows {
		prefs[i] = &models.UserPreferences{
			ID:                  rows[i].ID.Bytes,
			Device:              rows[i].Device,
			Theme:               models.ThemeFromDB(rows[i].Theme),
			OnboardingCompleted: rows[i].OnboardingCompleted,
			Language:            rows[i].Language,
			LastSyncAt:          rows[i].LastSyncAt.Time.UnixMilli(),
			CreatedAt:           rows[i].CreatedAt.Time.UnixMilli(),
			UpdatedAt:           rows[i].UpdatedAt.Time.UnixMilli(),
		}
	}
	return prefs, nil
}

func (r *UserPreferencesRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.UserPreferences, error) {
	row, err := r.queries.GetPreferencesByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.UserPreferences{
		ID:                  row.ID.Bytes,
		Device:              row.Device,
		Theme:               models.ThemeFromDB(row.Theme),
		OnboardingCompleted: row.OnboardingCompleted,
		Language:            row.Language,
		LastSyncAt:          row.LastSyncAt.Time.UnixMilli(),
		CreatedAt:           row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:           row.UpdatedAt.Time.UnixMilli(),
	}, nil
}

func (r *UserPreferencesRepo) Create(ctx context.Context, userID uuid.UUID, prefs *models.UserPreferences) error {
	return r.queries.CreateUserPreferences(ctx, db.CreateUserPreferencesParams{
		ID:                  pgtype.UUID{Bytes: prefs.ID, Valid: true},
		UserID:              pgtype.UUID{Bytes: userID, Valid: true},
		Device:              prefs.Device,
		Theme:               prefs.Theme.ToDB(),
		OnboardingCompleted: prefs.OnboardingCompleted,
		Language:            prefs.Language,
		CreatedAt:           pgtype.Timestamptz{Time: time.UnixMilli(prefs.CreatedAt), Valid: true},
		UpdatedAt:           pgtype.Timestamptz{Time: time.UnixMilli(prefs.UpdatedAt), Valid: true},
	})
}

func (r *UserPreferencesRepo) Update(ctx context.Context, prefs *models.UserPreferences) error {
	return r.queries.UpdateUserPreferences(ctx, db.UpdateUserPreferencesParams{
		Theme:     prefs.Theme.ToDB(),
		Language:  prefs.Language,
		UpdatedAt: pgtype.Timestamptz{Time: time.UnixMilli(prefs.UpdatedAt), Valid: true},
		ID:        pgtype.UUID{Bytes: prefs.ID, Valid: true},
	})
}

func (r *UserPreferencesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePreferences(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserPreferencesRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeletePreferencesByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
