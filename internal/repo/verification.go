package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationRepository struct {
	queries *db.Queries
}

func NewVerificationRepository(queries *db.Queries) *VerificationRepository {
	return &VerificationRepository{queries: queries}
}

func (r *VerificationRepository) Create(ctx context.Context, params *models.UserVerification) error {
	return r.queries.CreateVerification(ctx, db.CreateVerificationParams{
		UserID:    pgtype.UUID{Bytes: params.UserID, Valid: true},
		ActionID:  pgtype.UUID{Bytes: params.UserAction.ID},
		Token:     params.Token,
		ExpiresAt: pgtype.Timestamptz{Time: params.ExpiresAt, Valid: true},
	})
}

func (r *VerificationRepository) GetByToken(ctx context.Context, token string) (*models.UserVerification, error) {
	v, err := r.queries.GetVerificationByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !v.ActionID.Valid {
		return nil, errors.New("action id is not valid")
	}

	if !v.ActionName.Valid {
		return nil, errors.New("action name is not valid")
	}

	return &models.UserVerification{
		UserID: v.UserID.Bytes,
		Token:  v.Token,
		UserAction: &models.UserAction{
			ID:     v.ActionID.Bytes,
			Action: models.ActionFromString(v.ActionName.String),
		},
		ExpiresAt: v.ExpiresAt.Time,
		CreatedAt: v.CreatedAt.Time,
	}, nil
}

func (r *VerificationRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.queries.DeleteByToken(ctx, token)
}
