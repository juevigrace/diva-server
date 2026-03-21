package repo

import (
	"context"

	"github.com/juevigrace/diva-server/storage/db"
)

type UserPreferencesRepository struct {
	queries *db.Queries
}

func NewUserPreferencesRepository(queries *db.Queries) *UserPreferencesRepository {
	return &UserPreferencesRepository{queries: queries}
}

func (r *UserPreferencesRepository) Create(ctx context.Context, params *db.CreateUserPreferencesParams) error {

	return r.queries.CreateUserPreferences(ctx, *params)
}

func (r *UserPreferencesRepository) Update(ctx context.Context, params *db.UpdateUserPreferencesParams) error {
	return r.queries.UpdateUserPreferences(ctx, *params)
}
