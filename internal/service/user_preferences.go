package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPreferencesService struct {
	queries *db.Queries
}

func NewUserPreferencesService(queries *db.Queries) *UserPreferencesService {
	return &UserPreferencesService{
		queries: queries,
	}
}

func (s *UserPreferencesService) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPreferences, error) {
	rows, err := s.queries.GetPreferencesByUser(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		return nil, err
	}

	prefs := make([]*models.UserPreferences, len(rows))
	for i := range rows {
		prefs[i] = models.UserPrefsFromDB(&rows[i])
	}

	return prefs, nil
}

func (s *UserPreferencesService) GetByID(ctx context.Context, id uuid.UUID) (*models.UserPreferences, error) {
	row, err := s.queries.GetPreferencesByID(ctx, models.UUIDPtrToDB(&id))
	if err != nil {
		return nil, err
	}

	return models.UserPrefsFromDB(&row), nil
}

func (s *UserPreferencesService) Create(ctx context.Context, userID uuid.UUID, dto *dtos.CreateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:                  uuid.New(),
		Device:              dto.Device,
		Theme:               models.ThemeFromString(dto.Theme),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
	}

	return s.queries.CreateUserPreferences(ctx, *pref.DBCreate(userID))
}

func (s *UserPreferencesService) Update(ctx context.Context, id uuid.UUID, dto *dtos.UpdateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:       id,
		Theme:    models.ThemeFromString(dto.Theme),
		Language: dto.Language,
	}

	return s.queries.UpdateUserPreferences(ctx, *pref.DBUpdate())
}

func (s *UserPreferencesService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeletePreferences(ctx, models.UUIDPtrToDB(&id))
}

func (s *UserPreferencesService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.queries.DeletePreferencesByUser(ctx, models.UUIDPtrToDB(&userID))
}
