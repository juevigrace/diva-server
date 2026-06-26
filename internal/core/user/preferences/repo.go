package preferences

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core/user/permissions"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage"
)

type UserPreferencesRepo struct {
	store  storage.UserPreferenceStore
	upRepo *permissions.UserPermissionRepo
}

func NewUserPreferencesRepo(
	store storage.UserPreferenceStore,
	upRepo *permissions.UserPermissionRepo,
) *UserPreferencesRepo {
	return &UserPreferencesRepo{
		store:  store,
		upRepo: upRepo,
	}
}

func (s *UserPreferencesRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPreferences, error) {
	rows, err := s.store.GetPreferencesByUser(ctx, userID)
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
	row, err := s.store.GetPreferencesByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return models.UserPrefsFromDB(row), nil
}

func (s *UserPreferencesRepo) Create(ctx context.Context, session *models.Session, uid uuid.UUID, dto *dtos.CreateUserPreferencesDto) error {
	pref := &models.UserPreferences{
		ID:                  uuid.New(),
		Device:              dto.Device,
		Theme:               models.ThemeFromString(dto.Theme),
		OnboardingCompleted: dto.OnboardingCompleted,
		Language:            dto.Language,
	}
	if err := s.store.CreateUserPreferences(ctx, pref.DBCreate(uid)); err != nil {
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

	return s.store.UpdateUserPreferences(ctx, pref.DBUpdate())
}
