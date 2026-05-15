package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (r *UserPreferencesRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPreferences, error) {
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
			OnboardingCompleted: rows[i].Onboardingcompleted,
			Language:            rows[i].Language,
			LastSyncAt:          rows[i].Lastsyncat.Time.UnixMilli(),
			CreatedAt:           rows[i].Createdat.Time.UnixMilli(),
			UpdatedAt:           rows[i].Updatedat.Time.UnixMilli(),
		}
	}
	return prefs, nil
}

func (r *UserPreferencesRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserPreferences, error) {
	row, err := r.queries.GetPreferencesByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.UserPreferences{
		ID:                  row.ID.Bytes,
		Device:              row.Device,
		Theme:               models.ThemeFromDB(row.Theme),
		OnboardingCompleted: row.Onboardingcompleted,
		Language:            row.Language,
		LastSyncAt:          row.Lastsyncat.Time.UnixMilli(),
		CreatedAt:           row.Createdat.Time.UnixMilli(),
		UpdatedAt:           row.Updatedat.Time.UnixMilli(),
	}, nil
}

func (r *UserPreferencesRepository) GetByDevice(ctx context.Context, device string) (*models.UserPreferences, error) {
	row, err := r.queries.GetPreferencesByDevice(ctx, device)
	if err != nil {
		return nil, err
	}
	return &models.UserPreferences{
		ID:                  row.ID.Bytes,
		Device:              row.Device,
		Theme:               models.ThemeFromDB(row.Theme),
		OnboardingCompleted: row.Onboardingcompleted,
		Language:            row.Language,
		LastSyncAt:          row.Lastsyncat.Time.UnixMilli(),
		CreatedAt:           row.Createdat.Time.UnixMilli(),
		UpdatedAt:           row.Updatedat.Time.UnixMilli(),
	}, nil
}

func (r *UserPreferencesRepository) Create(ctx context.Context, userID uuid.UUID, prefs *models.UserPreferences) error {
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

func (r *UserPreferencesRepository) Update(ctx context.Context, prefs *models.UserPreferences) error {
	return r.queries.UpdateUserPreferences(ctx, db.UpdateUserPreferencesParams{
		Theme:     prefs.Theme.ToDB(),
		Language:  prefs.Language,
		UpdatedAt: pgtype.Timestamptz{Time: time.UnixMilli(prefs.UpdatedAt), Valid: true},
		ID:        pgtype.UUID{Bytes: prefs.ID, Valid: true},
	})
}

func (r *UserPreferencesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePreferences(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserPreferencesRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeletePreferencesByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
