package repo

import (
	"context"
	"errors"

	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationRepository struct {
	queries *db.Queries
}

func NewVerificationRepository(queries *db.Queries) *VerificationRepository {
	return &VerificationRepository{queries: queries}
}

func (r *VerificationRepository) Create(ctx context.Context, params *db.CreateVerificationParams) error {
	return r.queries.CreateVerification(ctx, *params)
}

func (r *VerificationRepository) GetByToken(ctx context.Context, token string) (*models.UserVerification, error) {
	v, err := r.queries.GetVerificationByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !v.ActionName.Valid {
		return nil, errors.New("action name is not valid")
	}

	return &models.UserVerification{
		UserID:    v.UserID.Bytes,
		Token:     v.Token,
		Action:    models.ActionFromString(v.ActionName.String),
		ExpiresAt: v.ExpiresAt.Time,
		CreatedAt: v.CreatedAt.Time,
	}, nil
}

func (r *VerificationRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.queries.DeleteByToken(ctx, token)
}
