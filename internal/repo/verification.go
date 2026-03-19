package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationRepository struct {
	queries *db.Queries
}

func NewVerificationRepository(queries *db.Queries) *VerificationRepository {
	return &VerificationRepository{queries: queries}
}

func (r *VerificationRepository) Create(ctx context.Context, userId *uuid.UUID) (*models.UserVerification, error) {
	token, err := util.GenerateOTPCode()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC()
	expiresAt := time.Now().Add(15 * time.Minute)

	verification := &models.UserVerification{
		UserID:    *userId,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	params := db.CreateParams{
		UserID:    pgtype.UUID{Bytes: verification.UserID, Valid: true},
		Token:     verification.Token,
		CreatedAt: pgtype.Timestamptz{Time: verification.CreatedAt, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: verification.ExpiresAt, Valid: true},
	}

	err = r.queries.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return verification, nil
}

func (r *VerificationRepository) GetByToken(ctx context.Context, token string) (*models.UserVerification, error) {
	v, err := r.queries.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.New("token not found")
	}
	return &models.UserVerification{
		UserID:    v.UserID.Bytes,
		Token:     v.Token,
		ExpiresAt: v.CreatedAt.Time,
		CreatedAt: v.ExpiresAt.Time,
	}, nil
}

func (r *VerificationRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.queries.DeleteByToken(ctx, token)
}
