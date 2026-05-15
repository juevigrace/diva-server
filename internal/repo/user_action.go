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

func (r *UserActionsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserAction, error) {
	row, err := r.queries.GetUserActionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.UserAction{
		ID:     row.ID.Bytes,
		Action: models.ActionFromString(row.Name),
	}, nil
}

func (r *UserActionsRepository) GetByUserAndName(ctx context.Context, userID uuid.UUID, action models.Action) (*models.UserAction, error) {
	row, err := r.queries.GetUserActionByUserAndName(ctx, db.GetUserActionByUserAndNameParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Name:   action.String(),
	})
	if err != nil {
		return nil, err
	}
	return &models.UserAction{
		ID:     row.ID.Bytes,
		Action: models.ActionFromString(row.Name),
	}, nil
}

func (r *UserActionsRepository) Create(ctx context.Context, userID uuid.UUID, ua *models.UserAction) error {
	return r.queries.CreateUserAction(ctx, db.CreateUserActionParams{
		ID:     pgtype.UUID{Bytes: ua.ID, Valid: true},
		Name:   ua.Action.String(),
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserActionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUserAction(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
