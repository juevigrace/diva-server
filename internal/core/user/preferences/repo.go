package preferences

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core/user/permissions"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPreferencesRepo struct {
	queries   *db.Queries
	upRepo *permissions.UserPermissionRepo
}

func NewUserPreferencesRepo(
	queries *db.Queries,
	upRepo *permissions.UserPermissionRepo,
) *UserPreferencesRepo {
	return &UserPreferencesRepo{
		queries:   queries,
		upRepo: upRepo,
	}
}

func (s *UserPreferencesRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPreferences, error) {
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

func (s *UserPreferencesRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.UserPreferences, error) {
	row, err := s.queries.GetPreferencesByID(ctx, models.UUIDPtrToDB(&id))
	if err != nil {
		return nil, err
	}

	return models.UserPrefsFromDB(&row), nil
}

func (s *UserPreferencesRepo) Create(ctx context.Context, session *models.Session, uid uuid.UUID, dto *dtos.CreateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:                  uuid.New(),
		Device:              dto.Device,
		Theme:               models.ThemeFromString(dto.Theme),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
	}
	if err := s.queries.CreateUserPreferences(ctx, *pref.DBCreate(uid)); err != nil {
		return err
	}
	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_PREFERENCES_WRITE]; ok {
			if err := s.upRepo.Delete(ctx, session.User.ID, perm.Permission.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *UserPreferencesRepo) Update(ctx context.Context, id uuid.UUID, dto *dtos.UpdateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:       id,
		Theme:    models.ThemeFromString(dto.Theme),
		Language: dto.Language,
	}

	return s.queries.UpdateUserPreferences(ctx, *pref.DBUpdate())
}
