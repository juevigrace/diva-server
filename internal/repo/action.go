package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type ActionRepository struct {
	queries *db.Queries
}

func NewActionRepository(queries *db.Queries) *ActionRepository {
	return &ActionRepository{queries: queries}
}

func (r *ActionRepository) Create(ctx context.Context, userAction *models.UserAction) error {
	params := db.CreateUserPendingActionParams{
		UserID:     pgtype.UUID{Bytes: userAction.UserID, Valid: true},
		ActionName: userAction.Action.String(),
	}
	return r.queries.CreateUserPendingAction(ctx, params)
}

func (r *ActionRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]db.DivaUserPendingAction, error) {
	return r.queries.GetUserPendingActions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (r *ActionRepository) Delete(ctx context.Context, userAction *models.UserAction) error {
	params := db.DeleteUserPendingActionParams{
		UserID:     pgtype.UUID{Bytes: userAction.UserID, Valid: true},
		ActionName: userAction.Action.String(),
	}
	return r.queries.DeleteUserPendingAction(ctx, params)
}
