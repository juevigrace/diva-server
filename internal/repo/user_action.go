package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserActionsRepository struct {
	queries *db.Queries
}

func NewUserActionsRepository(queries *db.Queries) *UserActionsRepository {
	return &UserActionsRepository{queries: queries}
}

func (r *UserActionsRepository) Create(ctx context.Context, userAction *models.UserAction) error {
	params := db.CreateUserPendingActionParams{
		UserID:     pgtype.UUID{Bytes: userAction.UserID, Valid: true},
		ActionName: userAction.Action.String(),
	}
	return r.queries.CreateUserPendingAction(ctx, params)
}

func (r *UserActionsRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]db.DivaUserPendingAction, error) {
	return r.queries.GetUserPendingActions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (r *UserActionsRepository) Delete(ctx context.Context, userAction *models.UserAction) error {
	params := db.DeleteUserPendingActionParams{
		UserID:     pgtype.UUID{Bytes: userAction.UserID, Valid: true},
		ActionName: userAction.Action.String(),
	}
	return r.queries.DeleteUserPendingAction(ctx, params)
}
