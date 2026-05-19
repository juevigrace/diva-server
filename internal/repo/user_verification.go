package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserVerificationRepo struct {
	queries *db.Queries
}

func NewUserVerificationRepo(queries *db.Queries) *UserVerificationRepo {
	return &UserVerificationRepo{queries: queries}
}

func (r *UserVerificationRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.UserActionVerification, error) {
	row, err := r.queries.GetVerification(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.UserActionVerification{
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
	}, nil
}

func (r *UserVerificationRepo) Create(ctx context.Context, v *models.UserActionVerification) error {
	return r.queries.CreateVerification(ctx, db.CreateVerificationParams{
		ActionID:  pgtype.UUID{Bytes: v.Action.ID, Valid: true},
		Token:     v.Token,
		ExpiresAt: pgtype.Timestamptz{Time: v.ExpiresAt, Valid: true},
	})
}

func (r *UserVerificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteVerification(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
