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
	params := db.CreateUserActionParams{
		UserID:     pgtype.UUID{Bytes: userAction.UserID, Valid: true},
		ActionName: userAction.Action.String(),
		ID:         pgtype.UUID{Bytes: userAction.ID, Valid: true},
	}
	return r.queries.CreateUserAction(ctx, params)
}

func (r *UserActionsRepository) DeleteByUser(ctx context.Context, userID *uuid.UUID) error {
	return r.queries.DeleteUserActions(ctx, pgtype.UUID{Bytes: *userID, Valid: true})
}

func (r *UserActionsRepository) CreateBatch(ctx context.Context, userActions []*models.UserAction) error {
	for _, ua := range userActions {
		if err := r.Create(ctx, ua); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserActionsRepository) Delete(ctx context.Context, userAction *models.UserAction) error {
	params := db.DeleteUserActionParams{
		UserID:     pgtype.UUID{Bytes: userAction.UserID, Valid: true},
		ActionName: userAction.Action.String(),
	}
	return r.queries.DeleteUserAction(ctx, params)
}

func (r *UserActionsRepository) GetAll(ctx context.Context, userID *uuid.UUID) ([]db.DivaUserPendingAction, error) {
	return r.queries.GetUserActions(ctx, pgtype.UUID{Bytes: *userID, Valid: true})
}

func (r *UserActionsRepository) GetOne(ctx context.Context, action models.Action, userID *uuid.UUID) (*models.UserAction, error) {
	params := db.GetUserActionParams{
		ActionName: action.String(),
		UserID:     pgtype.UUID{Bytes: *userID, Valid: true},
	}

	a, err := r.queries.GetUserAction(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.UserAction{
		ID:     a.ID.Bytes,
		Action: models.ActionFromString(a.ActionName),
		UserID: a.UserID.Bytes,
	}, nil
}
