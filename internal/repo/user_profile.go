package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserProfileRepo struct {
	queries *db.Queries
}

func NewUserProfileRepo(queries *db.Queries) *UserProfileRepo {
	return &UserProfileRepo{queries: queries}
}

func (r *UserProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	row, err := r.queries.GetUserProfileByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.UserProfile{
		FirstName: row.FirstName,
		LastName:  row.LastName,
		BirthDate: row.BirthDate.Time.UnixMilli(),
		Alias:     row.Alias,
		Bio:       row.Bio,
	}, nil
}

func (r *UserProfileRepo) Create(ctx context.Context, userID uuid.UUID, profile *models.UserProfile) error {
	return r.queries.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		BirthDate: pgtype.Timestamptz{Time: time.UnixMilli(profile.BirthDate), Valid: true},
		Alias:     profile.Alias,
		Bio:       profile.Bio,
	})
}

func (r *UserProfileRepo) Update(ctx context.Context, userID uuid.UUID, profile *models.UserProfile) error {
	return r.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		BirthDate: pgtype.Timestamptz{Time: time.UnixMilli(profile.BirthDate), Valid: true},
		Alias:     profile.Alias,
		Bio:       profile.Bio,
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserProfileRepo) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	return r.queries.UpdateUserProfileAvatar(ctx, db.UpdateUserProfileAvatarParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Avatar: avatar,
	})
}
