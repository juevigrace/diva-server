package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserActionsRepo struct {
	queries *db.Queries
}

func NewUserActionsRepo(queries *db.Queries) *UserActionsRepo {
	return &UserActionsRepo{queries: queries}
}

func (r *UserActionsRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserAction, error) {
	rows, err := r.queries.ListActionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	actions := make([]models.UserAction, len(rows))
	for i, row := range rows {
		actions[i] = models.UserAction{
			ID:     row.ID.Bytes,
			Name:   models.ActionFromString(row.Name),
			UserID: row.UserID.Bytes,
		}
	}

	return actions, nil
}

func (r *UserActionsRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.UserAction, error) {
	row, err := r.queries.GetUserActionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &models.UserAction{
		ID:     row.ID.Bytes,
		Name:   models.ActionFromString(row.Name),
		UserID: row.UserID.Bytes,
	}, nil
}

func (r *UserActionsRepo) GetByUserAndName(ctx context.Context, userID uuid.UUID, action models.Action) (*models.UserAction, error) {
	row, err := r.queries.GetUserActionByUserAndName(ctx, db.GetUserActionByUserAndNameParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Name:   action.String(),
	})
	if err != nil {
		return nil, err
	}

	return &models.UserAction{
		ID:     row.ID.Bytes,
		Name:   models.ActionFromString(row.Name),
		UserID: row.UserID.Bytes,
	}, nil
}

func (r *UserActionsRepo) Create(ctx context.Context, ua *models.UserAction) error {
	return r.queries.CreateUserAction(ctx, db.CreateUserActionParams{
		ID:     pgtype.UUID{Bytes: ua.ID, Valid: true},
		Name:   ua.Name.String(),
		UserID: pgtype.UUID{Bytes: ua.UserID, Valid: true},
	})
}

func (r *UserActionsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUserAction(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserActionsRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteUserActionByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
