package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.UpdateVerified(ctx, db.UpdateVerifiedParams{
		UserVerified: true,
		ID:           pgtype.UUID{Bytes: userID, Valid: true},
	})
}
