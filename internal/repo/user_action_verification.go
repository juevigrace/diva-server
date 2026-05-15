package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserActionVerificationRepository struct {
	queries *db.Queries
}

func NewVerificationRepository(queries *db.Queries) *UserActionVerificationRepository {
	return &UserActionVerificationRepository{queries: queries}
}

func (r *UserActionVerificationRepository) GetByActionAndToken(ctx context.Context, actionID uuid.UUID, token string) (*models.UserActionVerification, error) {
	row, err := r.queries.GetActionVerification(ctx, db.GetActionVerificationParams{
		ActionID: pgtype.UUID{Bytes: actionID, Valid: true},
		Token:    token,
	})
	if err != nil {
		return nil, err
	}
	return &models.UserActionVerification{
		Token:     row.Token,
		ExpiresAt: row.Expiresat.Time,
	}, nil
}

func (r *UserActionVerificationRepository) Create(ctx context.Context, actionID uuid.UUID, v *models.UserActionVerification) error {
	return r.queries.CreateVerification(ctx, db.CreateVerificationParams{
		ActionID:  pgtype.UUID{Bytes: actionID, Valid: true},
		Token:     v.Token,
		ExpiresAt: pgtype.Timestamptz{Time: v.ExpiresAt, Valid: true},
	})
}

func (r *UserActionVerificationRepository) Delete(ctx context.Context, actionID uuid.UUID, token string) error {
	return r.queries.DeleteVerification(ctx, db.DeleteVerificationParams{
		ActionID: pgtype.UUID{Bytes: actionID, Valid: true},
		Token:    token,
	})
}

func (r *UserActionVerificationRepository) DeleteByActionID(ctx context.Context, actionID uuid.UUID) error {
	return r.queries.DeleteVerificationByActionID(ctx, pgtype.UUID{Bytes: actionID, Valid: true})
}
